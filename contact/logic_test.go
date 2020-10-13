package emails

/* import (
	addresses "gitlab.com/digitalhouse-dev/core/notifications/email/postapp/addresses"
)

// ***************
// 	    Mocks
// ***************
var addrss = []addresses.Address{
	addresses.Address{
		ID:      1,
		Address: "nlcostamagna@gmail.com",
		Name:    "Nahuel",
		EmailID: 1,
		Type:    "TO",
	},
	addresses.Address{
		ID:      2,
		Address: "ncostamagna@digitalhouse.com",
		Name:    "Nahuel",
		EmailID: 1,
		Type:    "CC",
	},
}

var email1Ok = addressStruct{
	Name:    "Nahuel",
	Address: "nahuel@test.com",
}

var email2Ok = addressStruct{
	Name:    "Paco",
	Address: "paco@test.com",
}

var email1Err = addressStruct{
	Name:    "Nahuel",
	Address: "@test.com",
}

var email2Err = addressStruct{
	Name:    "Nahuel",
	Address: "",
}

func attachmentMocks() ([]attachStruct, int) {
	return []attachStruct{
		attachStruct{Content: "TGJASJKASASFDASAD", FileName: "Test.png", Type: "image/png"},
		attachStruct{Content: "GASKHLDUAUSIJ", FileName: "Test1.png", Type: "image/png"},
		attachStruct{Content: "ASADDSASADDA", FileName: "Test5.jpeg", Type: "image/jpeg"},
		attachStruct{Content: "ASKJJAHSLJALKSKJ", FileName: "Test6.jpeg", Type: "image/jpeg"},
	}, 4
}

func valuesEmailMocks() ([]addressStruct, []addressStruct, []addressStruct, addressStruct, addressStruct) {
	var bccOK = []addressStruct{
		addressStruct{Name: "Nahuel", Address: "testbcc@gmail.com"},
	}
	var ccOK = []addressStruct{
		addressStruct{Name: "Nahuel", Address: "testcc@gmail.com"},
	}
	var toOK = []addressStruct{
		addressStruct{Name: "Nahuel", Address: "testto1@gmail.com"},
		addressStruct{Name: "Nahuel", Address: "testto2@gmail.com"},
	}
	var fromOK = addressStruct{Name: "Nahuel", Address: "testfrom@gmail.com"}
	var replyToOK = addressStruct{Name: "Nahuel", Address: "testreplyTo@gmail.com"}

	return bccOK, ccOK, toOK, fromOK, replyToOK
} */

// ***************
// 	  Unit test
// ***************
/*
func TestCastAddressFunction(t *testing.T) {
	var emailToTest []*mail.Email

	assert.Nil(t, castAddress(&addrss, &emailToTest, email1Ok, "CC"))
	assert.NotNil(t, castAddress(&addrss, &emailToTest, email1Err, "TO"))
	assert.Nil(t, castAddress(&addrss, &emailToTest, email2Ok, "BCC"))
	assert.NotNil(t, castAddress(&addrss, &emailToTest, email2Err, "CC"))

	assert.EqualValues(t, len(addrss), 4)
	assert.EqualValues(t, len(emailToTest), 2)
} */

/* func TestCastEmailsFuncion(t *testing.T) {
	var email Email
	bcc, cc, to, from, replyTo := valuesEmailMocks()

	toSend, ccSend, bccSend, err := castEmails(&email, from, replyTo, to, cc, bcc)

	assert.EqualValues(t, len(toSend), 2)
	assert.EqualValues(t, len(ccSend), 1)
	assert.EqualValues(t, len(bccSend), 1)
	assert.Nil(t, err)
	assert.EqualValues(t, len(email.Addresses), 6)
} */
/*
func TestCastAttachmentsFuncion(t *testing.T) {
	var email Email
	mocks, sizeValue := attachmentMocks()

	attachments, err := castAttachments(&email, mocks)

	assert.EqualValues(t, len(attachments), sizeValue)
	assert.Nil(t, err)
} */

/*
 Necesitamos pensar las pruebas de integracion
func TestSendWithoutTo(t *testing.T) {
	bcc, cc, to, from, replyTo := valuesEmailMocks()
	Send
} */

// ***************
// 	 BenchMarks
// ***************
/*
func BenchmarkCastAddress(b *testing.B) {
	var emailToTest []*mail.Email
	for i := 0; i < b.N; i++ {
		_ = castAddress(&addrss, &emailToTest, email1Ok, "CC")
	}
}

func BenchmarkCastEmails(b *testing.B) {
	var email Email
	bcc, cc, to, from, replyTo := valuesEmailMocks()

	for i := 0; i < b.N; i++ {
		_, _, _, _ = castEmails(&email, from, replyTo, to, cc, bcc)
	}
}

func BenchmarkCastAttachments(b *testing.B) {
	var email Email
	mocks, _ := attachmentMocks()

	for i := 0; i < b.N; i++ {
		_, _ = castAttachments(&email, mocks)
	}
}
*/
