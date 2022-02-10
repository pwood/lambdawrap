# lambdawrap

[![license](https://img.shields.io/github/license/pwood/lambdawrap.svg)](https://github.com/pwood/lambdawrap/blob/master/LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)
[![Actions Status](https://github.com/pwood/lambdawrap/workflows/main/badge.svg)](https://github.com/pwood/lambdawrap/actions)

> A set of composable functions to reduce the repeated code required to implement AWS lambdas.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

Writing AWS Lambdas often require common tasks, such as unmarshalling SNS messages and then unmarshalling an inner SQS
message. This kind of code still requires testing in each individual lambda. This library attempts to remove that 
repeated complexity, as well as other useful utility functions. 

## Install

**This library makes use of Go generics, as such it is only compatible with Go >= 1.18!**

```shell
go get github.com/pwood/lambdawrap
```

## Usage

This library is under active development, it should not be used in production. 

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/pwood/lambdawrap/issues/new) or submit PRs.

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct/) Code
of Conduct.

## License

Copyright 2022 Peter Wood & Contributors

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "
AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.