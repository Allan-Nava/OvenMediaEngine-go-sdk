---
layout: default
title: Data Types
nav_order: 2
description: "Oven Media Engine Go SDK Official Documentation"
permalink: /data-types
last_modified_date: 2023-02-06T10:00:00+0000
---

# Data Types

## Primitives/Notations


| Type | Description | Examples |
| ------ | ------ | ------ |
| Short | 16bits integer | 12345 |
| Int | 32bits integer | 1234941932 |
| Long | 64bits integer | 391859818923919232311 |
| Float | 64bits real number | 3.5483 |
| String | A string | "Hello" |
| Bool | true/false | true |
| Timestamp (String)| A timestamp in ISO8601 format | "2021-01-01T11:00:00.000+09:00" |
| TimeInterval(Long) | A time interval (unit: milliseconds) | 349820 |
| IP address | "127.0.0.1" | "127.0.0.1" |
| RangedPort (String) | Port numbers with range (it can contain multiple ports and protocols)
start_port[-end_port][,start_port[-end_port][,start_port[-end_port]...]][/protocol] | "40000-40005/tcp"
"40000-40005"
"40000-40005,10000,20000/tcp" |
| Port (String) | A port number
start_port[/protocol]| "1935/tcp"
"1935"|



## Enum/Container Notations

#### Enum<Type> (String)

- An enum value named Type
- Examples
    - "value1"
    - "value2"


#### List<Type>
- An ordered list of Type
- Examples
    - [ Type, Type, ... ]

#### Map<KeyType, ValueType>
- An object consisting of Key-Value pairs
- Examples
    - { Key1: Value1, Key2: Value2 }