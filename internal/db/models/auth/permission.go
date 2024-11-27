package auth

type Permission string

const superAdmin Permission = "superadmin"

type Permissions interface {
    // AddPermissionForUser creates specified permission if it is not exists
    // and 
    AddPermissionForUser(userId int, perm Permission)     
    RemovePermissionForUser(userId int, perm Permission)
}
