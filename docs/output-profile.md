---
layout: default
title: Output Profile
nav_order: 5
description: "Oven Media Engine Go SDK Official Documentation"
permalink: /output-profile
last_modified_date: 2023-02-06T10:00:00+0000
---

# Output profile 

```
Creates OutputProfiles in the Application

Request Example:
POST http://1.2.3.4:8081/v1/vhosts/default/apps/app/outputProfiles

[
  {
    "name": "bypass_profile",
    "outputStreamName": "${OriginStreamName}",
    "encodes": {
      "videos": [
        {
          "bypass": true
        }
      ],
      "audios": [
        {
          "bypass": true
        }
      ]
    }
  }
]

```



````
Lists all output profile names in the Application

Request Example:
GET http://1.2.3.4:8081/v1/vhosts/default/apps/app/outputProfiles


````