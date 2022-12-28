---
title: "Namespaces"
date: 2021-08-31T17:26:20-04:00
draft: false
weight: 60
description: >
    Amazon Genomics CLI uses namespacing to prevent conflicts
---
Amazon Genomics CLI uses namespacing to prevent conflicts when there are multiple [users]( {{< relref "users" >}} ), [contexts]( {{< relref "contexts" >}} ), and [projects]( {{< relref "projects" >}} ) in the same AWS account and region.

In any given account and region, an individual user may have many projects with many deployed contexts all running at the
same time without conflict as long as:

1. No other user with the same Amazon Genomics CLI username exists in the same account and region.
2. All projects, used by that user, have a unique name.
3. All contexts within a project have a unique name.

## Shared Project Definitions

Project definitions can be shared between users. A simple way to achieve this is by putting the project YAML file and associated
workflow definitions into a source control system like Git. If two users in the same account and region start contexts
from the same project definition, these contexts are discrete and include the Amazon Genomics CLI username in the names of their respective
infrastructures.

Therefore, the following combination are allowed:

    userA -uses-> ProjectA -to-deploy-> ContextA
    userB -uses-> ProjectA -to-deploy-> ContextA

In the above example it is useful to think of these as two instances of Context A. Both share the same definition but the
instances do not have the same infrastructure.

## Tags

All Amazon Genomics CLI infrastructure is tagged with the `application-name` key and a value of `agc`
Aside from the core account infrastructure, all deployed infrastructure is tagged with the following key value pairs:


| Key            | Value                                                           |
| ---------------- | ----------------------------------------------------------------- |
| agc-project    | The name of the project in which the context is defined         |
| agc-user-id    | The unique username                                             |
| agc-user-email | The users email                                                 |
| agc-context    | The name of the context in which the infrastructure is deployed |
| agc-engine     | The name of the engine being run in the context                 |
| agc-engine-type| The workflow language run by the engine                         |

