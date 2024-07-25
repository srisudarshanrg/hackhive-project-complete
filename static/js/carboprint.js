document.addEventListener('DOMContentLoaded', function() {
  function handleSubmit(event) {
      event.preventDefault(); 

      // Retrieve all the variables with IDs
      var firstName = document.getElementById('firstName').value;
      var lastName = document.getElementById('lastName').value;
      var vehicleOwned = document.getElementById('vehicleOwned').value;
      var foodType = document.getElementById('foodType').value;
      var meatType = document.getElementById('meatType').value;
      var clothType = document.getElementById('clothType').value;
      var IntTravelPerYear = document.getElementById('IntTravelPerYear').value;
      var buildingType = document.getElementById('buildingType').value;
      var waterUsageDay = document.getElementById('waterUsageDay').value;
      var transportType = document.getElementById('transportType').value;
      var workCulture = document.getElementById('workCulture').value;
      var Gardens = document.getElementById('Gardens').value;
      var fuelTypeVehicle = document.getElementById('fuelTypeVehicle').value;
      var fuelTypeDomestic = document.getElementById('fuelTypeDomestic').value;
      var DailyTravel = document.getElementById('DailyTravel').value;

      // Process user's data
      var userTotalData = 0; // Initialize userTotalData to 0 from beginning

      if (vehicleOwned == "Yes") {
          userTotalData += 500;
      }

      // Add other calculations here...
      // Initial calculations
userTotalData += (foodType === "Vegetarian") ? 100 :
(foodType === "Lacto Vegetarian") ? 200 :
(foodType === "Pescetarian") ? 300 :
(foodType === "Flexitarian") ? 400 :
(foodType === "Non-Vegetarian") ? 500 : 0;

userTotalData += (meatType === "Beef") ? 800 :
(meatType === "Mutton") ? 750 :
(meatType === "Bacon") ? 700 :
(meatType === "Pork") ? 650 :
(meatType === "Turkey") ? 500 :
(meatType === "Duck") ? 450 :
(meatType === "Chicken") ? 400 :
(meatType === "Seafood") ? 300 : 0;

userTotalData += (clothType === "Silk") ? 700 :
(clothType === "Velvet") ? 660 :
(clothType === "Georgette") ? 640 :
(clothType === "Nylon") ? 580 :
(clothType === "Wool") ? 550 :
(clothType === "Rayon") ? 510 :
(clothType === "Denim") ? 470 :
(clothType === "Cotton") ? 440 : 0;

userTotalData += Number(IntTravelPerYear * 200);

userTotalData += (buildingType === "High-Rise") ? 500 :
(buildingType === "Low-Rise") ? 300 :
(buildingType === "Independent") ? 650 : 0;

userTotalData += Number(waterUsageDay);

userTotalData += (transportType === "Bus") ? 600 :
(transportType === "Bike") ? 400 :
(transportType === "Car") ? 500 :
(transportType === "Train") ? 750 : 0;

if (workCulture === "at home") {
userTotalData /= 2;
} else if (workCulture === "at office") {
userTotalData *= 2;
}

if (Gardens === "Yes") {
userTotalData /= 2;
} else if (Gardens === "No") {
userTotalData *= 1.5;
}

userTotalData += (fuelTypeVehicle === "Petrol") ? 500 :
(fuelTypeVehicle === "Diesel") ? 450 :
(fuelTypeVehicle === "Gasoline") ? 600 :
(fuelTypeVehicle === "Hydrogen") ? 250 : 0;

// Debugging - Output userTotalData
console.log("User Total Data: " + userTotalData);

      // You can follow the same structure as before for calculating userTotalData based on form inputs

      // Sample for logging all the data given
      console.log("FirstName: " + firstName + "\nLastName: " + lastName +
          "\nVehicle Owned: " + vehicleOwned + "\nFood Type: " + foodType + "\nMeat Type: " + 
          meatType + "\nCloth Type: " + clothType + "\nIntTravelPerYear: " + IntTravelPerYear + 
          "\nBuilding Type: " + buildingType + "\nWater Usage Per Day: " + waterUsageDay + "\nTransport Type: " + transportType +
          "\nGardens: " + Gardens + "\nFuel Type (Vehicle): " + fuelTypeVehicle + "\nWork Culture: " + workCulture + "\nFuel Type (Domestic): " + fuelTypeDomestic +
          "\nDaily Travel: " + DailyTravel);

      // Display userTotalData
      console.log("User Total Data: " + userTotalData);
      var cf = document.getElementById("cf")
      cf.innerText = userTotalData
  }
  // Add event listener to the form and handle the submit
  var form = document.getElementById('userForm');
  form.addEventListener('submit', handleSubmit);

});