# blist

Blist is an in-memory time series binary list package for Golang.

[![](https://img.shields.io/circleci/token/6b9dfa35f65f07aa855ecefffa1f3b05a803679f/project/abcum/blist/master.svg?style=flat-square)](https://circleci.com/gh/abcum/blist) [![](https://img.shields.io/badge/status-1.0.0-ff00bb.svg?style=flat-square)](https://github.com/abcum/blist) [![](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abcum/blist) [![](https://goreportcard.com/badge/github.com/abcum/blist?style=flat-square)](https://goreportcard.com/report/github.com/abcum/blist) [![](https://img.shields.io/badge/license-Apache_License_2.0-00bfff.svg?style=flat-square)](https://github.com/abcum/blist) 

#### Features

- In-memory binary list
- Store values by version number
- Delete values by version number
- Find the initial and the latest version
- Ability to insert items at any position in the list
- Find exact versions or seek to the closest version
- Select items by version number or retrieve latest value
- Sams efficiency as a btree when seeking for a specific version: O(log n) worst case
- Not as efficient as a tlist when majority of selects are for the initial or latest version: O(log n) worst case

#### Installation

```bash
go get github.com/abcum/blist
```
