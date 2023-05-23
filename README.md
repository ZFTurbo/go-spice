# go-spice

![Go](https://img.shields.io/badge/Go-1.16-blue.svg)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

go-spice is a program written in Go, designed to calculate the static (DC) IR-drop of integrated circuits. It serves as an alternative to C++ for performing these calculations efficiently and accurately.

## Features

- Calculates the static (DC) IR-drop of integrated circuits.
- Supports the Spice netlist format as the input data format.
- Generates output data in CSV format, with node names and corresponding values.

## Installation

Make sure you have Go 1.16 or higher installed. Then, run the following command to get go-spice:

```shell
git clone https://github.com/AlaieT/go-spice.git
```

## Usage

To use go-spice, you need to provide a Spice netlist file as the input. The program will then calculate the static IR-drop and generate a CSV file with the node names and their corresponding values.

```shell
go run ./cmd/pgsim -f test.spice -o ./ -p 1e-8
```

- `-f`: Path to the Spice netlist file.
- `-o`: Path to the output folder.
- `-p`: Modeling precision
- `-ms`: "Max modeling steps

The program will process the netlist file and generate the `[circuit_name].csv` file with the calculated IR-drop values.

## License

This project is licensed under the MIT License. See the [MIT LICENSE](./LICENSE) file for details.

## Contributing

Contributions to go-spice are welcome! If you find any bugs or have suggestions for improvements, please open an issue or submit a pull request.
