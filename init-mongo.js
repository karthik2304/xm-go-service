
// Switch to target database
db = db.getSiblingDB("companyDB");

// Create collections
db.createCollection("companies");
db.createCollection("events");
db.createCollection("users");


const companies = [
    {
      companyuuid: "550e8400-e29b-41d4-a716-446655440000",
      companyname: "TechCorp",
      description: "A technology company specializing in AI solutions.",
      totalemployees: 500,
      registered: true,
      type: "Corporations"
    },
    {
      companyuuid: "550e8400-e29b-41d4-a716-446655440001",
      companyname: "GreenNonProfit",
      description: "Nonprofit organization focused on sustainability.",
      totalemployees: 200,
      registered: true,
      type: "NonProfit"
    }
];
db.companies.insertMany(companies);
  
const usersdata = [
    {
        username: "karthik.coumar20@gmail.com",
        password: "$2a$10$MFSnaE7iWLx7datZq3Luh.XhQAX1Ryu4rXwP3St1aIXWBIlqJYi7q"
    },
    {
        username: "xmtest@gmail.com",
        password: "$2a$10$RdSETampBWuSRfAsLMQP6.7rSfqRKlKi7Vf/4xMCNbyeFk9CS/GFm"
    }
]
db.users.insertMany(usersdata)

print("✅ MongoDB: Collections created successfully!");
