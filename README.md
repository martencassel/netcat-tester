# Netstat-pid 

Netcat-tester tests network connections (TCP) against a file of dst_ip's and dst_ports.
The result are written to a output csv file.
 
## Usage

1. Test network connections specified from csv file. Output to a new csv file.

    ```sh
    cat test.csv
    dst_ip,dst_port
    8.8.8.8,53
    www.google.com,443

    ./netcat-tester -f test.csv -o output.csv
    ```

    ```sh
    cat test.csv
    
    ```

## Build


1. Fedora 22

    ```sh
    make build-fedora-32
    ```

2. Ubuntu 22.04

    ```sh
    make build-ubuntu-22.04
    ```
