package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta content="width=device-width, initial-scale=1" name="viewport">
        <title>Preflight CORS</title>
    </head>
    <body>
        <h1>Preflight CORS</h1>
        <div id=
                 "output"></div>
        <script>
            document.addEventListener('DOMContentLoaded', function () {
                fetch("http://localhost:4000/v1/tokens/authentication", {
                        method: "POST",
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            email: 'johnsmith1@gmail.com',
                            password: 'password'
                        })
                    }
                ).then(
                    function (response) {
                        response.text().then(function (text) {
                            document.getElementById("output").innerHTML = text;
                        });
                    },
                    function (err) {
                        document.getElementById("output").innerHTML = err;
                    }
                );
            });
        </script>
    </body>
</html>
`

func main() {
	addr := flag.String("addr", ":9000", "HTTP listen address")
	flag.Parse()

	log.Printf("Starting server on %s", *addr)

	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(html))
		if err != nil {
			log.Printf(err.Error())
		}
	}))
	log.Fatal(err)
}
