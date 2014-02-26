package main

import (
    "nowac/kegg"
    "fmt"
)

func main(){

    geneId := "10458"
    gene := kegg.GetGene(geneId);
    gene.Print()
   

    pathwayId := "hsa05200"
    pathway := kegg.GetPathway(pathwayId); 
    pathway.Print() 

    keggPathway := kegg.newKeggPathway(pathwayId)
    fmt.Println(keggPathway)   

}
