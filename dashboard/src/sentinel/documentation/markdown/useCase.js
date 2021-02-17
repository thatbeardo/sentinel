export const useCase = `

## Use Case

For this manual, we will setup all the Sentinel infrastructure keeping in mind this following use case. 
Danny is a Floor Manager at an e-commerce company. Danny has 2 employees reporting to him directly - Rusty, and Linus. Danny and his team use a portal to determine salaries. Given the hierarchy, as a manager, Danny should be able to edit salries and modify it. However, as an employee reporting to Danny, I should not be able to edit salaries. Danny and their team use MoneyManager – a payroll system developed in house. MoneyManager is written in React and the entire team logs into the same portal. When MoneyManager loads up, it consults Sentinel to get information about the currently logged in user and displays suitable UI components based on the role of who is logged in. For this example, we will create a bunch of resources, contexts and permissions to make sure Sentinel knows that Danny has privileges to edit salaries but others don't. 

`