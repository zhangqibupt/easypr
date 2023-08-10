## Introduction

`fwpr` is a command-line tool designed to simplify the process of creating Pull Requests and performing cherry-picking in Git repositories. It provides an efficient workflow to manage your code changes across branches, making the process of contribution and code synchronization smoother.

## Features

- Create Pull Requests from the current branch to specified target branches.
- Perform cherry-picking of commits across branches effortlessly.
- Simplify the management of code changes and synchronization.

## Installation
```sh
go install github.freewheel.com/qzhang/fwpr@latest
```

## Usage

#### Case 1: I want to create a pull request to merge `feature1` into `master` branch. Meanwhile cherry-pick to `V_6_57` and `V_6_57_1`
Just switch to `feature1` and run below command
```sh
fwpr create
```
It will let you choose the target branch(`master` in this case) and create a pull request to it. 

Then it will let you choose multiple cherry-pick branches (`V_6_57` and `V_6_57_1` in this case), then it will 
- checkout two new branches `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` from `V_6_57` and `V_6_57_1` respectively
- cherry-pick the commits from `feature1` to `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` respectively
- create two pull requests respectively

#### Case 2: The pull requests were created, but there are some new commits. I want to sync these new commits to all the pull requests.
Just switch to `feature1` and run below command
```sh
fwpr sync
```
It will let you choose the target branch(`master` in this case).
Then it will list the created cherry-pick branches and let you choose, `feature1_cp_to_V_6_57` and `feature1_cp_to_V_6_57_1` in this case. And then sync the new commits to them.
