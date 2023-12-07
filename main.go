package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	proto  = "tcp"
	dialer = &net.Dialer{Timeout: time.Second * 1}
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func Dial(dst_ip, dst_port string) (bool, error) {
	endpoint := net.JoinHostPort(dst_ip, dst_port)
	conn, err := dialer.Dial(proto, endpoint)
	if err != nil {
		return false, err
	}
	if conn != nil {
		defer conn.Close()
	}
	return true, nil
}

type ResultRecord struct {
	SrcIP   string
	DstIP   string
	DstPort string
	Success bool
}

func (rr *ResultRecord) String() string {
	return fmt.Sprintf("%s,%s,%s,%t", rr.SrcIP, rr.DstIP, rr.DstPort, rr.Success)
}

func main() {
	var verbose bool
	/*
		Parse command line arguments

		Usage:
		./netcat-tester -f <csv_file> -o <output_file> -v

	*/

	// Check if -v is set
	if len(os.Args) == 6 && os.Args[5] == "-v" {
		verbose = true
	} else {
		verbose = false
	}

	// Check if -f and -o are set
	if len(os.Args) != 5 {
		log.Fatalf("Usage: %s -f <csv_file> -o <output_file>", os.Args[0])
	}
	if os.Args[1] != "-f" || os.Args[3] != "-o" {
		log.Fatalf("Usage: %s -f <csv_file> -o <output_file>", os.Args[0])
	}
	// Parse csv_file and output_file
	csv_file := os.Args[2]
	output_file := os.Args[4]

	out_file, err := os.Create(output_file)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer out_file.Close()
	// Write header line
	fmt.Fprintf(out_file, "src_ip,dst_ip,dst_port,success\n")

	outboundIP := GetOutboundIP()

	// Read CSV file
	f, err := os.Open(csv_file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	// Read CSV header line
	r := csv.NewReader(f)
	_, err = r.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	for _, record := range records {
		// Print csv column
		dst_ip := record[0]
		dst_port := record[1]

		result := ResultRecord{
			SrcIP:   outboundIP.String(),
			DstIP:   dst_ip,
			DstPort: dst_port,
		}
		ok, err := Dial(dst_ip, dst_port)
		if err != nil {
			log.Warnf("Error dialing: %v", err)
		}
		if ok {
			result.Success = true
			fmt.Fprintf(out_file, "%s\n", result.String())
			if verbose {
				log.Infof("Success: %s %s %s", result.SrcIP, result.DstIP, result.DstPort)
			}
		} else {
			result.Success = false
			fmt.Fprintf(out_file, "%s\n", result.String())
			if verbose {
				log.Infof("Failure: %s %s %s", result.SrcIP, result.DstIP, result.DstPort)
			}
		}
	}
}
