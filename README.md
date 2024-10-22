# basePlatformSOMAS

<!-- Base platform for ICL SOMAS course -->

<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

basePlatformSOMAS is a generic 'base platform', in the form of a package, used to simplify the development of a self-organising, multi-agent system (SOMAS). This package provides a set of core 'building block' components:

- A 'base server' - an extensible server object for managing the environmental gamestate and regulating the internal state of...
- A 'base agent' - an extensible (multi-)agent object that can be injected into the server for simulation, and abstracts the core functions for agent-to-agent interactions through...
- A 'base message' - an extensible messaging component that helps to define a common language for multi-agent communication.

These building blocks come with 1. an interface definition for defining the core functions required for an operational multi-agent system and 2. the building block object itself, to provide a base implementation of these core functions, and allow extension.

For example, the `IServer` interface contains a set of functions that all servers in any multi-agent system must perform (the ability to add or remove agents, say), and the `BaseServer` object will give a base implementation of these methods.

It is then left to the user of this package to define an `ExtendedServer` which adds all of the relevant functionality for their multi-agent scenario, while composing the `BaseServer` to get access to the methods provided by this package.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

To begin working with basePlatformSOMAS, follow these simple example steps.

### Prerequisites

This package is entirely build in `GoLang (go)`. Your machine therefore needs a golang compiler installed. The relevant download link can be found [here](https://go.dev/doc/install).

[![Golang][go-shield]][go-url]

Since we have defined this repository as a package, it must be included in an existing `go` project. As such, you should be initialising a `go` repository in a folder of your choice with:

```sh
go mod init <your-module-name>
```

### Installation

With a working `go` repository, you can now include our package:

```sh
go get github.com/MattSScott/basePlatformSOMAS/v2
```

and, for peace of mind, tidy the imports:

```sh
go mod tidy
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

You will now have access to the modules provided by this package:

- The _server_ module can be imported into a file with: `github.com/MattSScott/basePlatformSOMAS/v2/pkg/server`
- The _agent_ module can be imported into a file with: `github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent`
- The _message_ module can be imported into a file with: `github.com/MattSScott/basePlatformSOMAS/v2/pkg/message`

For more examples, and a **much, much more detailed write-up** of the package, please refer to the [User Manual](https://github.com/MattSScott/basePlatformSOMAS/blob/main/basePlatformSOMASv2.0.pdf) or, if you're using an outdated version of the package, refer to the [Past Manuals](https://github.com/MattSScott/basePlatformSOMAS/tree/main/Past%20Manuals).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

basePlatformSOMAS is an open source project made for SOMAS programmers, by SOMAS programmers. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this package better, please fork the repo and create a pull request. If you would rather leave the implementation of this suggestion to us, you can also simply open an issue with the tag "enhancement".

To make a pull request in a helpful way:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

[Matthew Scott](https://profiles.imperial.ac.uk/matthew.scott18) - matthew.scott18@imperial.ac.uk

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

This work was produced by the following developers:

- [Matthew Scott](https://github.com/MattSScott)
- [Ana Dimoska](https://github.com/ADimoska)
- [Mikayel Suvaryan](https://github.com/mika111)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/MattSScott/basePlatformSOMAS.svg
[contributors-url]: https://github.com/MattSScott/basePlatformSOMAS/graphs/contributors
[issues-shield]: https://img.shields.io/github/issues/MattSScott/basePlatformSOMAS.svg?color=orange
[issues-url]: https://github.com/MattSScott/basePlatformSOMAS/graphs/issues
[license-shield]: https://img.shields.io/github/license/MattSScott/basePlatformSOMAS.svg
[license-url]: https://github.com/MattSScott/basePlatformSOMAS/blob/main/LICENSE.txt
[go-shield]: https://img.shields.io/badge/GoLang-blue?logo=go
[go-url]: https://go.dev
