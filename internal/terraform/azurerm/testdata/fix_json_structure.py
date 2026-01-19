#!/usr/bin/env python3
"""
Script to fix JSON test data files to match the golden file structure.
Ensures that only leaf nodes appear in costComponents, and nested components
are properly organized under subResources.
"""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Any, Tuple, Set


def parse_golden_file(golden_path: str) -> Dict[str, Any]:
    """
    Parse a golden file to extract the hierarchical structure.
    Returns a dict mapping resource names to their expected structure.
    Structure includes ordered list of subresources.
    """
    if not os.path.exists(golden_path):
        return {}
    
    with open(golden_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    resources = {}
    current_resource = None
    parent_stack = []  # Stack of (level, name, index) tuples
    
    for line in content.split('\n'):
        # Skip empty lines and header lines
        if not line.strip() or 'Monthly Qty' in line or line.strip().startswith('Name'):
            continue
        
        # Check if this is a resource name (starts with azurerm_ and no tree chars)
        if 'azurerm_' in line and not re.match(r'^ *[├└│]', line):
            resource_name = line.strip().split()[0]
            current_resource = resource_name
            resources[current_resource] = {
                'costComponents': [],
                'subResources': []  # Changed to list to preserve order
            }
            parent_stack = []
            continue
        
        if not current_resource:
            continue
        
        # Match tree structure lines
        tree_match = re.match(r'^( +)(?:│\s+)?(├─|└─)\s+(.+)$', line)
        if not tree_match:
            continue
        
        prefix = tree_match.group(1)
        content_part = tree_match.group(3).strip()
        
        # Determine depth by the full prefix including continuation chars
        full_prefix = line[:line.index('├─') if '├─' in line else line.index('└─')]
        
        depth = 0
        if len(full_prefix) >= 4:
            depth = 1
        elif len(full_prefix) > 1:
            depth = (len(full_prefix) - 1) // 3
        
        # Extract the name
        if 'Monthly cost depends' in content_part:
            name = content_part.split('Monthly cost depends')[0].strip()
        elif '  ' in content_part:
            name = content_part.split('  ')[0].strip()
        else:
            parts = re.split(r'\s{2,}', content_part)
            name = parts[0].strip() if parts else content_part.strip()
        
        # Determine if this is a leaf (has quantity/cost info) or a parent (subResource)
        is_leaf = bool(re.search(r'\d[\d,]*\s+(hours|GB|emails|events|notifications|messages|calls|10k operations|months|vCPU)', content_part)) or 'Monthly cost depends' in content_part
        
        # Update parent stack
        parent_stack = [p for p in parent_stack if p[0] < depth]
        
        if depth == 0:
            # Top-level item
            if is_leaf:
                resources[current_resource]['costComponents'].append(name)
            else:
                # This is a subResource - add as new entry
                sr_index = len(resources[current_resource]['subResources'])
                resources[current_resource]['subResources'].append({
                    'name': name,
                    'costComponents': []
                })
                parent_stack.append((depth, name, sr_index))
        else:
            # Child item - should be under the most recent parent
            if parent_stack:
                parent_name, parent_idx = parent_stack[-1][1], parent_stack[-1][2]
                if is_leaf:
                    resources[current_resource]['subResources'][parent_idx]['costComponents'].append(name)
    
    return resources


def get_all_cost_component_names(obj: Any) -> Set[str]:
    """Recursively get all cost component names from a resource or subresource."""
    names = set()
    
    if isinstance(obj, dict):
        if 'costComponents' in obj:
            for cc in obj.get('costComponents', []):
                if isinstance(cc, dict) and 'name' in cc:
                    names.add(cc['name'])
        
        if 'subResources' in obj:
            for sr in obj.get('subResources', []):
                names.update(get_all_cost_component_names(sr))
    
    return names


def build_structure_from_golden(golden_structure: Dict[str, Any]) -> Tuple[List[str], List[Dict[str, Any]]]:
    """
    Build lists of expected cost components and subresources from golden structure.
    Returns (top_level_cost_components, ordered_subresources).
    """
    top_level_ccs = golden_structure.get('costComponents', [])
    subresources = golden_structure.get('subResources', [])
    
    return top_level_ccs, subresources


def fix_resource_structure(resource: Dict[str, Any], golden_structure: Dict[str, Any]) -> Dict[str, Any]:
    """
    Fix a single resource's structure based on the golden structure.
    """
    if not golden_structure:
        return resource
    
    top_level_ccs, expected_subresources = build_structure_from_golden(golden_structure)
    
    # Collect all cost components from the resource (flattened)
    all_components = []
    for cc in resource.get('costComponents', []):
        all_components.append(cc)
    
    for sr in resource.get('subResources', []):
        for cc in sr.get('costComponents', []):
            all_components.append(cc)
    
    # Build new structure based on golden structure
    new_cost_components = []
    new_subresources = []
    
    # Track which components we've already used
    component_queue = list(all_components)
    
    # First, assign top-level cost components
    for expected_cc_name in top_level_ccs:
        # Find matching component in queue
        for i, cc in enumerate(component_queue):
            if cc.get('name') == expected_cc_name:
                new_cost_components.append(cc)
                component_queue.pop(i)
                break
    
    # Then, assign subresources in order
    for expected_sr in expected_subresources:
        sr_name = expected_sr['name']
        sr_cc_names = expected_sr['costComponents']
        
        sr_cost_components = []
        for expected_cc_name in sr_cc_names:
            # Find matching component in queue
            for i, cc in enumerate(component_queue):
                if cc.get('name') == expected_cc_name:
                    sr_cost_components.append(cc)
                    component_queue.pop(i)
                    break
        
        if sr_cost_components:
            new_subresources.append({
                'name': sr_name,
                'costComponents': sr_cost_components,
                'subResources': []
            })
    
    return {
        'name': resource['name'],
        'costComponents': new_cost_components,
        'subResources': new_subresources
    }


def fix_json_file(json_path: str, golden_path: str) -> Tuple[bool, str]:
    """
    Fix a JSON file based on its golden file.
    Returns (changed, message).
    """
    # Parse golden file
    golden_structures = parse_golden_file(golden_path)
    
    if not golden_structures:
        return False, f"Could not parse golden file or no structures found"
    
    # Load JSON file
    with open(json_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    if 'resources' not in data:
        return False, "No resources found in JSON"
    
    # Track if we made changes
    changed = False
    resources_fixed = 0
    
    # Fix each resource
    new_resources = []
    for resource in data['resources']:
        resource_name = resource.get('name', '')
        
        if resource_name in golden_structures:
            fixed_resource = fix_resource_structure(resource, golden_structures[resource_name])
            
            # Check if it changed
            if json.dumps(resource, sort_keys=True) != json.dumps(fixed_resource, sort_keys=True):
                changed = True
                resources_fixed += 1
            
            new_resources.append(fixed_resource)
        else:
            new_resources.append(resource)
    
    if changed:
        # Write back
        data['resources'] = new_resources
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2)
        
        return True, f"Fixed {resources_fixed} resources"
    
    return False, "No changes needed"


def main():
    """Main function to process all test data files."""
    testdata_dir = Path(__file__).parent
    
    # Find all golden files
    golden_files = list(testdata_dir.glob('**/*.golden'))
    
    print(f"Found {len(golden_files)} golden files")
    print("=" * 80)
    
    total_fixed = 0
    total_skipped = 0
    total_errors = 0
    
    for golden_path in sorted(golden_files):
        # Find corresponding JSON file(s)
        base_name = golden_path.stem
        json_path = golden_path.with_suffix('.json')
        
        if not json_path.exists():
            continue
        
        relative_path = json_path.relative_to(testdata_dir)
        
        try:
            changed, message = fix_json_file(str(json_path), str(golden_path))
            
            if changed:
                print(f"✓ {relative_path}: {message}")
                total_fixed += 1
            else:
                print(f"  {relative_path}: {message}")
                total_skipped += 1
        
        except Exception as e:
            print(f"✗ {relative_path}: Error - {str(e)}")
            total_errors += 1
    
    print("=" * 80)
    print(f"Summary:")
    print(f"  Fixed: {total_fixed}")
    print(f"  Skipped: {total_skipped}")
    print(f"  Errors: {total_errors}")


if __name__ == '__main__':
    main()
