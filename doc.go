// Package zebra is a web framework for Go. It is designed to be simple and easy to use.
//
// To get started, create a new Zebra instance and serve it:
//
//	func main() {
//			app, err := zebra.New()
//			if err != nil {
//				log.Fatalln(err)
//			}
//
//			log.Fatalln(app.ListenAndServe(":8080"))
//	}
package zebra
