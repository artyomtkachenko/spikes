package main

import (
	"fmt"

	"gopkg.in/rana/ora.v3"
)

func main() {
	env, srv, ses, err := ora.NewEnvSrvSes("user/pass@host:1521/service-name", nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to DB.......")
		s, err := ses.Prep("select 1 from dual", ora.U8)
		defer s.Close()
		if err != nil {
			panic(err)
		}
		rset, err := s.Qry()
		if err != nil {
			panic(err)
		}

		row := rset.NextRow()
		if row != nil {
			fmt.Printf("Got respones %d\n", row[0])
		}
		if rset.Err != nil {
			panic(rset.Err)
		}

	}
	defer env.Close()
	defer srv.Close()
	defer ses.Close()
}
