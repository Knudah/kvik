package main

import (
    "log"
    "os"
    "encoding/csv"
    "io"
    "strconv"

)

type Dataset struct {
    Bg Background
    Exprs Expression
}

type Background struct {
    IdInfo map[string] Info    
}

type Info struct {
    Lab int
    Id string
    Ctrl string
    Et  int
    En int
    Stage int
}

type Expression struct {
    Genes []string
    IdExpression map[string] [] float64
    GeneExpression map[string] [] float64
}

func NewDataset( path string ) Dataset {

    exprsFilename := path + "/exprs.csv"
    bgFilename := path + "/background.csv"

    bg, err := generateBackgroundDataset(bgFilename)
    if err != nil {
        log.Print(err)
    }
    
    exprs, err := generateExpressionDataset(exprsFilename)
    if err != nil {
        log.Print(err) 
    }
    
    ds := Dataset{bg, exprs}

    return ds
    
}

func generateExpressionDataset(filename string) (Expression, error) { 

    var IdExpression map[string][]float64
    var GeneExpression map[string][]float64

    exprs := Expression{}

    exprsfile, err := os.Open(filename)
    if err != nil {
        return exprs, err
    }  
    defer exprsfile.Close() 

    reader := csv.NewReader(exprsfile) 
    firstRow := true

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        if firstRow {
            exprs.Genes = record
            // the lengths here are maybe a bit off?
            IdExpression = make(map[string][]float64, len(record)-1)
            GeneExpression = make(map[string][]float64, len(record)-1)
            firstRow = false 
        } else { 

            id := record[0]
            expression := toFloats(record[1:])
            
            // store an id to expression mapping. 
            IdExpression[id] = expression
            
            // Store the expression value for this specific gene and id
            // combination
            for i, _ := range(expression){
                gene := exprs.Genes[i]
                GeneExpression[gene] = append(GeneExpression[gene],
                expression[i])
            }
        }
    }

    exprs.IdExpression = IdExpression
    exprs.GeneExpression = GeneExpression

    return exprs, nil
}

func (ds Dataset) PrintDebugInfo() {
    exprs := ds.Exprs

    log.Print("Generated dataset with ", len(exprs.IdExpression["900229_1"]),
                " genes and ",len(exprs.GeneExpression[exprs.Genes[0]]),
                " case/ctrl pairs")


}


func generateBackgroundDataset(filename string) (Background, error) { 
    
    bg := Background{}

    bgfile, err := os.Open(filename)
    if err != nil {
        return bg, err
    }
    defer bgfile.Close() 

    
    reader := csv.NewReader(bgfile) 
    firstRow := true

    var idinfo map[string] Info

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        if firstRow {
            
            idinfo = make(map[string] Info, len(record)) 

            firstRow = false 
        } else { 
            var lab, et, en, stage int

            lab, err := strconv.Atoi(record[0])
            if err != nil {
                log.Print("Error parsing ") 
                lab = 0
            }

            id := record[1]
            casectrl := record[2]
            et, err = strconv.Atoi(record[3])
            if err != nil {
                et = -1
            }

            en, err = strconv.Atoi(record[4])
            if err != nil {
                en = -1
            }

            stage, err = strconv.Atoi(record[5])
            if err != nil {
                stage = -1
            }
            
            info := Info{
                Lab: lab,
                Id: id,
                Ctrl: casectrl,
                Et: et,
                En: en,
                Stage: stage,
            }
            idinfo[id] = info

        }
    }

    bg.IdInfo = idinfo


    return bg, nil
    
    
}



func toFloats(input []string) []float64{
    output := make([]float64, len(input))

    var err error

    for i, _ := range(input){
        output[i], err = strconv.ParseFloat(input[i], 32)
        if err != nil {
            log.Panic("Parsing of float went bad: ", err)
        }
    }

    return output

}



