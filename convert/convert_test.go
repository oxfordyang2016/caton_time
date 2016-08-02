package main

import (
    "testing"  
)





type Testinstuct struct {
    inputjson          string
    expectedOutput     string
}


//----------------------------------------------------------------------------------------------------------------------->

var passvar = []Testinstuct{
    {`{ "from":"a", "cmd": "Login", "type": "rsp", "seq": 4294967295, "token": "EF02JLGFA09GVNG21F" }`,"EF02JLGFA09GVNG21F"},
 
}


func TestConver(t *testing.T) {
    for _, test := range passvar {  // 3.4
        if output := Convert(test.inputjson); output != test.expectedOutput {//test func 
            // 4
            t.Errorf("Expected SayTwice to return (%s), Found (%s)\n", test.expectedOutput, output)
        }
    }
}
