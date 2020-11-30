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
    "body": "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en-GB\"><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\" /><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"/><title>Time to Visit</title><style>p {color: blueviolet;font-weight: bolder;}</style></head><body style=\"margin: 0; padding: 0;\"><table align=\"center\" border=\"1\" cellpadding=\"0\" cellspacing=\"0\" width=\"600\" style=\"border-collapse: collapse;\"><tr><td style=\"padding: 40px 0 30px 0;\"><p style=\"margin: 0;\">Hi {{ Name }},</p><br><p><b>How are you doing !!</b><br><I>Nice connecting with you.</I> I'm scheduling a tech talk this month.</p><p>And feel delighted to have you as our guest speaker. I know your competence and proud of your antecedence :) <br>Thank you as we look forward to your favourable response.</p><br><br><p>Thanks, <br>Harish</p></td><td bgcolor=\"#ee4c50\" style=\"padding: 30px 30px;\"><p style=\"margin: 0;\">Tech Talk Invitation</p></td></tr></table></body></html>"
}


