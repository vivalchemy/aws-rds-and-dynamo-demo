Here are the `curl` commands for each API endpoint:
### **1. Get All Plants**
```sh
curl -X GET http://localhost:8080/plants
```

### **2. Get a Specific Plant**
```sh
curl -X GET http://localhost:8080/plants/{id}
```
Replace `{id}` with the actual plant ID.

### **3. Create a New Plant**
```sh
curl -X POST http://localhost:8080/plants \
     -H "Content-Type: application/json" \
     -d '{
          "name": "Aloe Vera",
          "scientific_name": "Aloe barbadensis miller",
          "family": "Asphodelaceae",
          "type": "Succulent",
          "sunlight_required": "Full Sun",
          "water_interval": 7,
          "height": 0.5,
          "native": "Africa",
          "indoor": true
         }'
```

### **4. Update an Existing Plant**
```sh
curl -X PUT http://localhost:8080/plants/{id} \
     -H "Content-Type: application/json" \
     -d '{
          "name": "Updated Aloe Vera",
          "scientific_name": "Aloe barbadensis miller",
          "family": "Asphodelaceae",
          "type": "Succulent",
          "sunlight_required": "Partial Sun",
          "water_interval": 5,
          "height": 0.6,
          "native": "Africa",
          "indoor": false
         }'
```
Replace `{id}` with the actual plant ID.

### **5. Delete a Plant**
```sh
curl -X DELETE http://localhost:8080/plants/{id}
```
Replace `{id}` with the actual plant ID.

Let me know if you need modifications! 🚀
