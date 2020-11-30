// Payload Sample for various commands: -JSON payload:

// Payload for AddClient command
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

// Payload for Send Mail or Schedule Mail command
[{
"templateId":2,
"clientId":7
},{
"templateId":2,
"clientId":8
}]

// Payload Add template or Update template commamd
{
    "name": "TemplateName",
    "subject": "Lets start campaigning",
    "body": "Hi {{ Name }}, How are you doing !"
}


