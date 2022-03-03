# strongDM Go SDK Examples

This is the examples repository for the [strongDM Go SDK](https://github.com/strongdm/strongdm-sdk-go).

---
> **NOTE:**  
> To increase flexibility when managing a large volume of Resources, Role Grants have
been deprecated in favor of Access Rules, which allow you to grant access based
on Resource Tags and Type.
>
> Previously, you would grant a Role access to specific resources by ID via Role
Grants. Now, when using Access Rules, the best practice is to give Roles access to Resources based on 
type and tags.
>
>The following examples demonstrate Dynamic Access Rules with tags and resource types, as well as Static Access Rules. If it is _necessary_ to grant access to specific Resources in the same way as Role Grants did, you can use Resource IDs directly in Access Rules.

---
