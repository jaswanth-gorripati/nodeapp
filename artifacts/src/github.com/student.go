package main


import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type StudentChainCode struct {
}


type Education struct{
	Degree	string `json:"degree"`
	Board string `json:"board"`
	Institute string  `json:"Institute`
	YearOfPassout int  `json:"yearOfPassout"`
	Score	float64 `json:"score"`
	AddedBy string `json:"addedBy"`
	AddedTime time.Time `json:"addedTime"`
	WhoCanupdate WhoCanUpdate `josn:"whoCanUpdate"`
}
type WhoCanUpdate struct{
	Name string `json:"userName"`
	WorkingPlace string `json:"worksAt"`
	WorkingAs string `json:"worksAs"`
}

type Student struct{
	CreatedBy string `json:"createdBy"`
	CreatedTime time.Time `json:"createdTime"`
	ProfilePic string `json:"profilePic"`
	Name string `json:"username"`
	DateOfBirth time.Time `json:"dateOfBirth`
	Gender string `json:"gender"`
	Education []Education `json:"education"`
}

//---------------------------------------------------------------------------------------------
// main
// --------------------------------------------------------------------------------------------
func main() {
	err := shim.Start(new(StudentChainCode))
	if err != nil {
		fmt.Printf("Error while starting Student Chaincode - %s", err)
	}
}

func(t *StudentChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Student's chaincode is starting up")
	fmt.Println(" - ready for action");
	id, err := cid.GetID(stub);
	fmt.Printf("id obtained %v",id);
	mspid, err := cid.GetMSPID(stub);
	fmt.Printf("mspid obtained %v",mspid);
	val, ok, err := cid.GetAttributeValue(stub, "position");
	if err != nil {
		fmt.Println("There was an error trying to retrieve the attribute");
	}else if !ok {
		fmt.Printf("The client identity does not possess the attribute %v",val);
	}else{
		fmt.Printf("value obtained %v",val);
	}
	return shim.Success(nil)
}

func (t *StudentChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" Invoke ");
	id, err := cid.GetID(stub);
	fmt.Printf("id obtained %v \n",id);
	mspid, err := cid.GetMSPID(stub);
	fmt.Printf("mspid obtained %v \n",mspid);
	val, ok, err := cid.GetAttributeValue(stub, "position");
	if err != nil {
		fmt.Println("There was an error trying to retrieve the attribute");
	}else if !ok {
		fmt.Println("The client identity does not possess the attribute");
	}else{
		fmt.Printf("value obtained %v",val);
	}
	fmt.Printf("value obtained %v",val);
	function, args := stub.GetFunctionAndParameters()
	if function == "init" {
		return t.Init(stub)	
	}else if function == "register" {
		return register(stub, args);
	}else if function == "getHistory" {
		return getHistory(stub, args)
	}else if function == "addEducation" {
		return addEducation(stub , args)
	}else if function == "search" {
		return search(stub,args)	
	}else if function == "getStudentDetails" {
		return getStudentDetails(stub,args)	
	}else if function == "getDetails" {
		return getDetails(stub,args)	
	}else if function == "cwcu" {
		return createWhoCanUpdate(stub,args)	
	}else if function == "updateEdu" {
		return updateEducationDetails(stub,args)	
	}
	
	fmt.Println("Received unknown invoke function name -" + function)
	return shim.Error("Received unknown invoke function name -'" + function + "'")
}
func (t *StudentChainCode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" QUERY ");
	function, args := stub.GetFunctionAndParameters()
	if function == "getStudentDeatils" {
		return getStudentDetails(stub,args)	
	}else if function == "getDetails" {
		return getDetails(stub,args)	
	}
	fmt.Println("Received unknown query function name -" + function)
	return shim.Error("Received unknown query function name -'" + function + "'")
}

// ==============================Register Student =======================================

func register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Registration of student in process... ")
	var education []Education
	var degree, board, institute string
	if len(args) != 14 && len(args) !=7{
		return shim.Error("Incorrect number of arguments. Expecting 9. studentId  and Students details to set")
	}

	// input sanitation
	//err = sanitize_arguments(args)
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	studentId :=args[0]
	displayPicture :=args[1]
	name :=args[2]                            
	dob,err := time.Parse("2006-01-02", args[3])
	if err != nil {
		return shim.Error("2st argument must be a date string and should be in 2006-01-02 format")
	}
	gender :=strings.ToLower(args[4])
	if gender!="male" && gender!="female"{
		return shim.Error("Gender must be male or female in 3rd argument")	
	}
	createdBy := args[5]
	createdTime ,err :=time.Parse("2006-01-02", args[6])
	if err != nil {
		return shim.Error("2st argument must be a date string and should be in 2006-01-02 format")
	}
	if len(args)==14 {
		degree =args[7]
		board = args[8]
		institute = args[9]
		passOut, err := strconv.Atoi(args[10])
		if err != nil {
			return shim.Error("7st argument must be a numeric string i,e. Year Of PassOut")
		} 
		score, err := strconv.ParseFloat(args[11], 64)
		if err != nil {
		   	return shim.Error("cannot convert to float ")
		}
		addedBy := args[12]
		addedTime,err :=time.Parse("2006-01-02", args[13])
		updateBy := WhoCanUpdate{"",institute,"admin"}
		if err != nil {
			return shim.Error("2st argument must be a date string and should be in 2006-01-02 format")
		}
		fmt.Println("edu details :")
		fmt.Println(addedBy)
		fmt.Println(addedTime)
		temp := Education{degree,board,institute,passOut,score,addedBy,addedTime,updateBy}
		education = append(education,temp)
		fmt.Println(education)
	}
	
	
	// To check if student already exists
	studentDetailsAsBytes, err := stub.GetState(studentId)
	if err != nil {
		fmt.Println("Registration Failed")
		return shim.Error("Failed to get student details: " + err.Error())
	} else if studentDetailsAsBytes != nil {
		fmt.Println("This student Id  already exists: "+args[0])
		fmt.Println("Registration Failed")
		return shim.Error("This student Id already exists: " + args[0])
	}
	//Assinging to student json
	//education := Education{degree,board,institute,passOut,score}
	student := Student{createdBy,createdTime,displayPicture,name,dob,gender,education}
	studentDetailsJSONasBytes, err := json.Marshal(student)
	if err != nil {
		fmt.Println("error while Json marshal")
		return shim.Error(err.Error())
	}
	//write the variable into the ledger
	err = stub.PutState(studentId, studentDetailsJSONasBytes)         
	if err != nil {
		return shim.Error(err.Error())
	}
	// CREATING A COMPOSITE KEY
	indexName := "gender~name~degree~board~institute"
	searchStudentsIndexKey, err := stub.CreateCompositeKey(indexName, []string{student.Gender,student.Name,degree,board,institute})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(searchStudentsIndexKey, value)
	var str []string
	str = append(str,"1")
	//getHistory(stub,str)
	fmt.Println("Registration Successfull ---> Record created " )
	return shim.Success(nil)
}

//==================================== Add education =====================================

func addEducation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	studentId := args[0]
	var updateBy WhoCanUpdate
	fmt.Println("Request received to Add Education deatils to the student")
	studentDetailsAsBytes, err := stub.GetState(studentId)
	if err != nil {
		fmt.Println("Updating Education Failed ")
		return shim.Error("Failed to get student details: " + err.Error())
	} else if studentDetailsAsBytes == nil {
		fmt.Println(" Student Details not found "+args[0])
		fmt.Println(" Updating Education Failed ")
		return shim.Error("This student Id already exists: " + args[0])
	}
	if (len(args)%7) != 1 {
		return shim.Error(" Incorrect number of arguments. Education details Missing some information")
	}
	studentDetailsAsJson := Student{}
	json.Unmarshal([]byte(studentDetailsAsBytes),&studentDetailsAsJson)
	var temp  Education
	var educationDetails []Education
	initialDetails :=  studentDetailsAsJson.Education
	for i:=0; i<len(initialDetails);i++{
		educationDetails = append(educationDetails,initialDetails[i])
	}

	for i :=1; i< len(args) ;i+=7{
		yop,err := strconv.Atoi(args[i+3])
		if err != nil {
			return shim.Error("argument must be a numeric string i,e. Year Of PassOut")
		} 
		degreeScore, err := strconv.ParseFloat(args[i+4], 64)
		if err != nil {
		   	return shim.Error("cannot convert to float ")
		}
		addedBy := args[i+5]
		addedTime, err := time.Parse("2006-01-02", args[i+6])
		updateBy = WhoCanUpdate{"",args[i+2],"admin"}
		temp = Education{args[i],args[i+1],args[i+2],yop,degreeScore,addedBy,addedTime,updateBy}
		educationDetails = append(educationDetails,temp)	
	}
	student := Student{studentDetailsAsJson.CreatedBy ,studentDetailsAsJson.CreatedTime,studentDetailsAsJson.ProfilePic,studentDetailsAsJson.Name,studentDetailsAsJson.DateOfBirth,studentDetailsAsJson.Gender,educationDetails}
	studentDetailsJSONasBytes, err := json.Marshal(student)
	if err != nil {
		fmt.Println("error while Json marshal")
		return shim.Error(err.Error())
	}
	//write the variable into the ledger
	err = stub.PutState(studentId, studentDetailsJSONasBytes)         
	if err != nil {
		return shim.Error(err.Error())
	}
	
	fmt.Println(educationDetails)
	fmt.Println("Student Education Details Updated...  :)")
	return shim.Success(nil)
}

//========================== Search by gender ============================================

func search (stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1. i,e Gender")
	}
	searchByGender := args[0]
	studentDetailsSearchIterator, err := stub.GetStateByPartialCompositeKey("gender~name~degree~board~institute", []string{searchByGender})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer studentDetailsSearchIterator.Close()

	// Iterate through result set and for each student details found
	var i int
	for i = 0; studentDetailsSearchIterator.HasNext(); i++ {
		
		responseRange, err := studentDetailsSearchIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedGender := compositeKeyParts[0]
		returnedName := compositeKeyParts[1]
		returnedDegre := compositeKeyParts[2]
		returnedBoard:= compositeKeyParts[3]
		returnedInstitute := compositeKeyParts[4]
		fmt.Printf(" %v ---> gender:%s Name:%s Degree:%s Board: %s institute:%s \n", i, returnedGender, returnedName,returnedDegre,returnedBoard,returnedInstitute)
	}

	fmt.Printf("- end of search  for %s Total results : %v \n",searchByGender,i)
	return shim.Success(nil)
}

//==============================Get History ==================================

func getHistory(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Student Id")
	}

	studentId := args[0]
	resultsIterator, err := stub.GetHistoryForKey(studentId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		historicValue,err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(historicValue.TxId)
		buffer.WriteString("\"")
		buffer.WriteString(", \"StudentDetails\":")

		// historicValue is a JSON , so we write as-is
		if historicValue != nil {
			buffer.WriteString(string(historicValue.Value))
		} else {
			buffer.WriteString("null")
		}
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
		
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryOfStudent returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//================getStudentDeatils (only works in couchDB)===================

func  getStudentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Requesting To display Student details...")
	//   0
	// "queryString"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse,err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- Found Student details : \n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func getDetails (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creator,err := stub.GetCreator()
	if err!=nil{
		return shim.Error(err.Error())
	}
	fmt.Println("creator : "+string(creator))
	/*signedProposal,err := stub.GetSignedProposal()
	if err!=nil{
		return shim.Error(err.Error())
	}
	fmt.Println(signedProposal)*/
	return shim.Success(nil)
}
func createWhoCanUpdate (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5{
		return shim.Error("Incorrect number of arguments. Expecting 5. who can update details")
	}
	studentId := args[0]
	studentDetailsAsBytes, err := stub.GetState(studentId)
	if err != nil {
		fmt.Println("Updating authorization Failed ")
		return shim.Error("Failed to get student details: " + err.Error())
	} else if studentDetailsAsBytes == nil {
		fmt.Println(" Updating Authorization Failed ")
		return shim.Error("This student Id already exists: ")
	}
	studentDetailsAsJson := Student{}
	forEduDetails, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2st argument must be a numeric string i,e. ")
	} 
	json.Unmarshal([]byte(studentDetailsAsBytes),&studentDetailsAsJson)
	var educationDetails []Education
	//whoCanUpdate := WhoCanUpdate{args[2],args[3],args[4]}
	initialDetails :=  studentDetailsAsJson.Education
	for i:=0; i<len(initialDetails);i++{
		educationDetails = append(educationDetails,initialDetails[i])
	}
	educationDetails[forEduDetails].WhoCanupdate = WhoCanUpdate{args[2],args[3],args[4]}
	student := Student{studentDetailsAsJson.CreatedBy ,studentDetailsAsJson.CreatedTime,studentDetailsAsJson.ProfilePic,studentDetailsAsJson.Name,studentDetailsAsJson.DateOfBirth,studentDetailsAsJson.Gender,educationDetails}
	studentDetailsJSONasBytes, err := json.Marshal(student)
	if err != nil {
		fmt.Println("error while Json marshal")
		return shim.Error(err.Error())
	}
	//write the variable into the ledger
	err = stub.PutState(studentId, studentDetailsJSONasBytes)         
	if err != nil {
		return shim.Error(err.Error())
	}
	
	fmt.Println(educationDetails)
	fmt.Println("Student Education Details Updated...  :)")
	return shim.Success(nil)

}
func updateEducationDetails (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 8 {
		return shim.Error(" Incorrect number of arguments. Details Missing ")
	}
	studentId := args[0]
	forEduDetails, err := strconv.Atoi(args[1])
	updatingBy := args[2]
	fmt.Println("Request received to Update Education details of the student")
	studentDetailsAsBytes, err := stub.GetState(studentId)
	if err != nil {
		fmt.Println("Updating Education Failed ")
		return shim.Error("Failed to get student details: " + err.Error())
	} else if studentDetailsAsBytes == nil {
		fmt.Println(" Student Details not found "+args[0])
		fmt.Println(" Updating Education Failed ")
		return shim.Error("This student Id already exists: " + args[0])
	}
	studentDetailsAsJson := Student{}
	json.Unmarshal([]byte(studentDetailsAsBytes),&studentDetailsAsJson)
	var educationDetails []Education
	initialDetails :=  studentDetailsAsJson.Education
	for i:=0; i<len(initialDetails);i++{
		educationDetails = append(educationDetails,initialDetails[i])
	}
	fmt.Println(educationDetails[forEduDetails].WhoCanupdate.Name)
	if educationDetails[forEduDetails].WhoCanupdate.Name != strings.ToUpper(updatingBy){
		fmt.Println("Updating Education Failed ")
		return shim.Error("Anuthorized request: "+updatingBy+" cannot perform update")	
	}else {
		yop,err := strconv.Atoi(args[6])
		if err != nil {
			return shim.Error("argument must be a numeric string i,e. Year Of PassOut")
		} 
		degreeScore, err := strconv.ParseFloat(args[7], 64)
		if err != nil {
		   	return shim.Error("cannot convert to float ")
		}
		addedBy := educationDetails[forEduDetails].AddedBy
		addedTime := educationDetails[forEduDetails].AddedTime
		updateBy := WhoCanUpdate{educationDetails[forEduDetails].WhoCanupdate.Name,educationDetails[forEduDetails].WhoCanupdate.WorkingPlace,educationDetails[forEduDetails].WhoCanupdate.WorkingAs}
		educationDetails[forEduDetails] = Education{args[3],args[4],args[5],yop,degreeScore,addedBy,addedTime,updateBy}
		//educationDetails = append(educationDetails,temp)	
	}
	
	student := Student{studentDetailsAsJson.CreatedBy ,studentDetailsAsJson.CreatedTime,studentDetailsAsJson.ProfilePic,studentDetailsAsJson.Name,studentDetailsAsJson.DateOfBirth,studentDetailsAsJson.Gender,educationDetails}
	studentDetailsJSONasBytes, err := json.Marshal(student)
	if err != nil {
		fmt.Println("error while Json marshal")
		return shim.Error(err.Error())
	}
	//write the variable into the ledger
	err = stub.PutState(studentId, studentDetailsJSONasBytes)         
	if err != nil {
		return shim.Error(err.Error())
	}
	
	fmt.Println(educationDetails)
	fmt.Println("Student Education Details Updated by :")
	fmt.Println(updatingBy)
	return shim.Success(nil)
}
