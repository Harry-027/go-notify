// Payload Sample for various commands: -JSON payload:

// AddClient command
[{
    "name": "Jon",
    "mailID": "Jon@gmail.com",
    "phone": 910088998899,
    "preference": "daily"
}, {
    "name": "Dove",
    "mailID": "Dove@gmail.com",
    "phone": 910088998899,
    "preference": "weekly"
}]

// Send Mail or Schedule Mail command
[{
"templateId":2,
"clientId":7
},{
"templateId":2,
"clientId":8
}]

// Add template or Update template commamd
{
    "name": "TemplateName",
    "subject": "Lets start campaigning",
    "body": "Hi {{ Name }}, How are you doing !"
}


