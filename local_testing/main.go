package main

//
//import (
//	"cloudru-containerapps-mcp/internal/application"
//	"cloudru-containerapps-mcp/internal/domain"
//	"log"
//)
//
//func main() {
//	ca := application.NewContainerAppsApplication()
//
//	// cas, err := ca.GetListContainerApps(
//	// 	"a9e46dcd-b00a-4a87-8ec2-028b31931c7b",
//	// 	domain.Credentials{
//	// 		KeyID:     "464450b9fa5f9e4ae8fe961d21ffacc9",
//	// 		KeySecret: "f5d8b6eb3df36232a1985678d9edf7a7",
//	// 	})
//	// log.Print(cas)
//	// if err != nil {
//	// 	log.Fatalf(err.Error())
//	// }
//	cas_, err := ca.GetContainerApp(
//		"a9e46dcd-b00a-4a87-8ec2-028b31931c7b",
//		"manage-photos-django-app",
//		domain.Credentials{
//			KeyID:     "464450b9fa5f9e4ae8fe961d21ffacc9",
//			KeySecret: "f5d8b6eb3df36232a1985678d9edf7a7",
//		})
//	log.Print(cas_)
//	if err != nil {
//		log.Fatalf(err.Error())
//	}
//	// cas, err := ca.GetContainerApp(
//	// 	"a9e46dcd-b00a-4a87-8ec2-028b31931c7b",
//	// 	"",
//	// 	domain.Credentials{
//	// 		KeyID:     "464450b9fa5f9e4ae8fe961d21ffacc9",
//	// 		KeySecret: "f5d8b6eb3df36232a1985678d9edf7a7",
//	// 	})
//	// log.Print(cas)
//	// if err != nil {
//	// 	log.Fatalf(err.Error())
//	// }
//}
