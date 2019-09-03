# eloquactl
Command line tool for Eloqua administration, data extraction and imports

# Examples

```bash
# exports Eloqua activities of type EmailSend
eloquactl export activities \
  --type=es \
  --since=2019-01-01
  
# exports Eloqua contacts
eloquactl export contacts \
  --fields='EmailAddress:{{Contact.Field(C_EmailAddress)}}'
  
# export Custom Data Object by its name
eloquactl export cdo mycdo1

# export Custom Data Object by its id
eloquactl export cdo 15
```
