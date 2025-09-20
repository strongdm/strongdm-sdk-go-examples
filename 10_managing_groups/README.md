# Managing Groups

This directory contains comprehensive CRUD examples for managing groups and their relationships with accounts and roles using the StrongDM SDK.

## Examples

### Groups CRUD
[`groups_crud`](./groups_crud) - Complete CRUD operations for Groups
- **Create**: Create new groups
- **Read**: List and filter groups
- **Update**: Modify group properties
- **Delete**: Remove groups
- Includes resource cleanup after demonstration

### AccountsGroups CRUD  
[`accounts_groups_crud`](./accounts_groups_crud) - Complete CRUD operations for AccountsGroups
- **Create**: Link accounts (users) to groups
- **Read**: List and filter account-group relationships
- **Delete**: Remove account-group relationships
- Creates prerequisite accounts and groups
- Includes complete resource cleanup

### GroupsRoles CRUD
[`groups_roles_crud`](./groups_roles_crud) - Complete CRUD operations for GroupsRoles
- **Create**: Link groups to roles
- **Read**: List and filter group-role relationships by group or role
- **Delete**: Remove group-role relationships
- Creates prerequisite groups and roles
- Includes complete resource cleanup

## Prerequisites

All examples require:
- StrongDM API keys set as environment variables:
  - `SDM_API_ACCESS_KEY`
  - `SDM_API_SECRET_KEY`
- Go 1.21 or later

## Usage

Each example can be run independently and demonstrates the complete lifecycle:

```bash
# Groups CRUD example
cd groups_crud
go run main.go

# AccountsGroups CRUD example
cd accounts_groups_crud  
go run main.go

# GroupsRoles CRUD example
cd groups_roles_crud
go run main.go
```

## Features

- **Complete CRUD Operations**: Each example demonstrates Create, Read, Update (where applicable), and Delete operations
- **Comprehensive Listing**: Shows how to list all resources and filter by specific criteria
- **Resource Management**: All examples create test resources, demonstrate operations, and clean up afterwards
- **Real-world Usage**: Examples show practical patterns for managing groups in production environments
- **Error Handling**: Proper error handling and logging throughout
- **Resource Dependencies**: Examples that require prerequisite resources (accounts, groups, roles) create them automatically