## Introduction

`fwpr` is a command-line tool designed to simplify the process of creating Pull Requests and performing cherry-picking in Git repositories. It provides an efficient workflow to manage your code changes across branches, making the process of contribution and code synchronization smoother.

## Demo
Below demo shows how to create pull requests to merge `feature1` into `master` branch. Meanwhile, cherry-pick to `V_6_57` and `V_6_57_1`

https://media.github.freewheel.tv/user/347/files/2d38fa75-3cd2-4eac-bdfd-9512375664a7

## Installation
```sh
go install github.freewheel.com/qzhang/fwpr@latest
```

## Usage
```
 # generate Pull Requests to multiple branches
fwpr create

# sync the new commits to all the created Pull Requests
fwpr sync

# set default assignees for Pull Requests
fwpr config set-assignees <ldap1> <ldap2> 
```
### Case 1:
> I want to create a pull request to merge `feature1` into `master` branch. Meanwhile cherry-pick to `V_6_57` and `V_6_57_1`

Just switch to `feature1` and run
```sh
fwpr create
```
It will let you choose the target branch(`master` in this case) and create a pull request to it. 

Then it will let you choose multiple cherry-pick branches (`V_6_57` and `V_6_57_1` in this case), then it will 
- Checkout two new branches `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` from `V_6_57` and `V_6_57_1` respectively
- Cherry-pick the commits from `feature1` to `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` respectively
- Create two pull requests respectively

### Case 2: 
> The pull requests were created, but there are some new commits. I want to sync these new commits to all the pull requests.

Just switch to `feature1` and run below command
```sh
fwpr sync
```
It will let you choose the target branch(`master` in this case).
Then it will list the created cherry-pick branches and let you choose, `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` in this case. And then sync the new commits to them.

## Limitations
- Currently, when cherry-picking, if there are conflicts, you need to resolve them manually.