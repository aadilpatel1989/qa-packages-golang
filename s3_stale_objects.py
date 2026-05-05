#!/usr/bin/env python3
"""
Script to find S3 objects that:
1. Were NOT accessed for more than 30 days
2. Were created before January 2026

Output is dumped into a CSV file.

Prerequisites:
- pip install boto3
- AWS credentials configured (via ~/.aws/credentials, env vars, or IAM role)
- S3 Storage Lens or S3 Inventory must be enabled for "last accessed" data,
  OR you can use the object's LastModified as a proxy.

Note: AWS S3 does not natively track "last accessed time" unless you enable
      S3 Storage Lens access metrics or use S3 Server Access Logging.
      This script uses LastModified as a proxy for last access time.
      If you have S3 Inventory with "Last Access Date" enabled, see the
      alternative approach in the comments below.
"""

import boto3
import csv
import os
from datetime import datetime, timezone, timedelta
import argparse

# ======================== CONFIGURATION ========================
# These will be set via command-line arguments
BUCKET_NAME = None
PREFIX = None
OUTPUT_CSV = None
DAYS_THRESHOLD = None
CREATED_BEFORE = None
# ===============================================================

def get_s3_client():
    """Initialize and return an S3 client."""
    return boto3.client('s3')

def get_all_objects(s3_client, bucket_name, prefix=""):
    """
    Retrieve all objects from the specified S3 bucket using pagination.
    
    Args:
        s3_client: boto3 S3 client
        bucket_name: Name of the S3 bucket
        prefix: Optional prefix to filter objects
        
    Returns:
        List of object metadata dictionaries
    """
    all_objects = []
    paginator = s3_client.get_paginator('list_objects_v2')
    
    page_config = {
        'Bucket': bucket_name,
    }
    if prefix:
        page_config['Prefix'] = prefix

    print(f"Fetching objects from bucket: {bucket_name}...")
    
    for page in paginator.paginate(**page_config):
        if 'Contents' in page:
            all_objects.extend(page['Contents'])
    
    print(f"Total objects found: {len(all_objects)}")
    return all_objects

def get_object_details(s3_client, bucket_name, key):
    """
    Get detailed metadata for a specific S3 object using HeadObject.
    This can provide additional metadata if needed.
    
    Args:
        s3_client: boto3 S3 client
        bucket_name: Name of the S3 bucket
        key: Object key
        
    Returns:
        Object metadata dictionary
    """
    try:
        response = s3_client.head_object(Bucket=bucket_name, Key=key)
        return response
    except Exception as e:
        print(f"Error fetching details for {key}: {e}")
        return None

def filter_stale_objects(objects, days_threshold, created_before):
    """
    Filter objects that:
    1. Have not been accessed (LastModified) for more than `days_threshold` days
    2. Were created before `created_before` date
    
    Args:
        objects: List of S3 object metadata
        days_threshold: Number of days threshold for staleness
        created_before: datetime object - only include objects created before this date
        
    Returns:
        List of stale objects with relevant details
    """
    stale_objects = []
    now = datetime.now(timezone.utc)
    threshold_date = now - timedelta(days=days_threshold)
    
    print(f"\nFiltering objects:")
    print(f"  - Not accessed for more than {days_threshold} days (before {threshold_date.strftime('%Y-%m-%d')})")
    print(f"  - Created before {created_before.strftime('%Y-%m-%d')}")
    
    for obj in objects:
        key = obj['Key']
        last_modified = obj['LastModified']
        size = obj['Size']
        
        # Skip "folder" markers (zero-byte objects ending with /)
        if key.endswith('/') and size == 0:
            continue
        
        # Check conditions:
        # 1. Last accessed (using LastModified as proxy) > 30 days ago
        # 2. Created before January 2026
        is_stale = last_modified < threshold_date
        is_created_before_cutoff = last_modified < created_before
        
        if is_stale and is_created_before_cutoff:
            days_since_access = (now - last_modified).days
            stale_objects.append({
                'Key': key,
                'LastModified': last_modified.strftime('%Y-%m-%d %H:%M:%S UTC'),
                'Size_Bytes': size,
                'Size_MB': round(size / (1024 * 1024), 4),
                'Days_Since_Last_Access': days_since_access,
                'StorageClass': obj.get('StorageClass', 'STANDARD'),
                'ETag': obj.get('ETag', '').strip('"'),
            })
    
    print(f"  - Stale objects found: {len(stale_objects)}")
    return stale_objects

def export_to_csv(stale_objects, output_file):
    """
    Export the list of stale objects to a CSV file.
    
    Args:
        stale_objects: List of dictionaries containing stale object info
        output_file: Path to the output CSV file
    """
    if not stale_objects:
        print("\nNo stale objects found. CSV file will not be created.")
        return
    
    fieldnames = [
        'Key',
        'LastModified',
        'Size_Bytes',
        'Size_MB',
        'Days_Since_Last_Access',
        'StorageClass',
        'ETag'
    ]
    
    with open(output_file, mode='w', newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(stale_objects)
    
    file_size = os.path.getsize(output_file)
    print(f"\n✅ CSV file created successfully!")
    print(f"   File: {output_file}")
    print(f"   Size: {file_size} bytes")
    print(f"   Records: {len(stale_objects)}")

def print_summary(stale_objects):
    """Print a summary of the stale objects found."""
    if not stale_objects:
        return
    
    total_size_bytes = sum(obj['Size_Bytes'] for obj in stale_objects)
    total_size_gb = total_size_bytes / (1024 ** 3)
    max_days = max(obj['Days_Since_Last_Access'] for obj in stale_objects)
    avg_days = sum(obj['Days_Since_Last_Access'] for obj in stale_objects) / len(stale_objects)
    
    # Group by storage class
    storage_classes = {}
    for obj in stale_objects:
        sc = obj['StorageClass']
        storage_classes[sc] = storage_classes.get(sc, 0) + 1
    
    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)
    print(f"  Total stale objects:        {len(stale_objects)}")
    print(f"  Total size:                 {total_size_gb:.4f} GB ({total_size_bytes:,} bytes)")
    print(f"  Max days since access:      {max_days} days")
    print(f"  Avg days since access:      {avg_days:.1f} days")
    print(f"  Storage class breakdown:")
    for sc, count in storage_classes.items():
        print(f"    - {sc}: {count} objects")
    print("=" * 60)

def parse_date(date_str):
    """
    Parse a date string in YYYY-MM-DD format to a datetime object.
    
    Args:
        date_str: String in YYYY-MM-DD format
        
    Returns:
        datetime object
    """
    try:
        return datetime.strptime(date_str, '%Y-%m-%d').replace(tzinfo=timezone.utc)
    except ValueError:
        raise ValueError("Date must be in YYYY-MM-DD format")

def main():
    """Main function to orchestrate the S3 stale object detection."""
    parser = argparse.ArgumentParser(description='Find S3 objects not accessed for more than X days and created before a specific date.')
    parser.add_argument('--bucket', required=True, help='S3 bucket name')
    parser.add_argument('--prefix', default="", help='Optional prefix to filter objects (folder path)')
    parser.add_argument('--output', default='s3_stale_objects.csv', help='Output CSV file name')
    parser.add_argument('--days', type=int, required=True, help='Number of days threshold for staleness')
    parser.add_argument('--created-before', required=True, help='Date objects were created before (YYYY-MM-DD)')

    args = parser.parse_args()

    # Set global configuration from command-line arguments
    BUCKET_NAME = args.bucket
    PREFIX = args.prefix
    OUTPUT_CSV = args.output
    DAYS_THRESHOLD = args.days
    CREATED_BEFORE = parse_date(args.created_before)

    print("=" * 60)
    print("S3 STALE OBJECT FINDER")
    print("=" * 60)
    print(f"Bucket:          {BUCKET_NAME}")
    print(f"Prefix:          {PREFIX if PREFIX else '(none - scanning entire bucket)'}")
    print(f"Days Threshold:  {DAYS_THRESHOLD}")
    print(f"Created Before:  {CREATED_BEFORE.strftime('%Y-%m-%d')}")
    print(f"Output File:     {OUTPUT_CSV}")
    print("=" * 60)
    
    # Initialize S3 client
    s3_client = get_s3_client()
    
    # Get all objects from the bucket
    all_objects = get_all_objects(s3_client, BUCKET_NAME, PREFIX)
    
    if not all_objects:
        print("No objects found in the bucket. Exiting.")
        return
    
    # Filter stale objects
    stale_objects = filter_stale_objects(all_objects, DAYS_THRESHOLD, CREATED_BEFORE)
    
    # Export to CSV
    export_to_csv(stale_objects, OUTPUT_CSV)
    
    # Print summary
    print_summary(stale_objects)

# ======================== ALTERNATIVE APPROACH ========================
# If you have S3 Inventory with "Last Access Date" enabled (requires
# S3 Storage Lens or Intelligent-Tiering access patterns), you can use:
#
# def get_last_access_from_inventory(s3_client, bucket_name, key):
#     """
#     Get last access time from S3 Object metadata using
#     x-amz-object-attributes (requires opt-in).
#     """
#     response = s3_client.get_object_attributes(
#         Bucket=bucket_name,
#         Key=key,
#         ObjectAttributes=['ObjectParts', 'StorageClass', 'ObjectSize']
#     )
#     return response
# =====================================================================


if __name__ == "__main__":
    main()