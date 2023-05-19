const BASE_URL = 'https://pluvio-api.herokuapp.com/api/v1';

window.addEventListener("load", (event) => {

    $('#Burkina_Faso').on('click', async function(){ 
        const $select = document.querySelector('#regions');
        $select.value = 'Burkina Faso';
        await rainRegion("Burkina Faso");
    });
    
    $('#Mali').on('click',async function(){ 
        const $select = document.querySelector('#regions');
        $select.value = 'Mali';
        await rainRegion("Mali");
    });

    $('#button').on('click', async function () {
        const phoneNumber = document.querySelector('#fmphone').value;
        const amount = document.querySelector('#fmrain').value;
        const location = document.querySelector('#regions').value;

        const response = await reportData(
            phoneNumber,
            location,
            parseInt(amount)
        );

        alert(response.message);
    });

});

function changeColour(elem) {
    var country_id = elem.id
    var colour = "#004400";
    //with selector:
    var country = document.getElementById(country_id);
    country.style.fill = colour;
    
    //or simply with elem arg passed into event handler:
    //elem.style.fill = colour;
}

async function rainRegion(value) {
    document.getElementById('rain-region').innerHTML = value;
    const dayData = await loadData('day', value);
    const weekData = await loadData('week', value);

    document.getElementById('rain-yesterday').innerHTML = dayData.rainfall;
    document.getElementById('rain-lastweek').innerHTML = weekData.rainfall;
}

async function loadData(timeRange, location) {
    if (location.includes('Faso')) {
        location = 'Faso';
    }
    const url = `${BASE_URL}/rain/${timeRange}?location=${location}`;

    const fetchOptions = {
        headers: {
            "Accept": "application/json"
        }
    };

    const response = await fetch(url, fetchOptions);
    const data = await response.json();
    return data;
}

async function reportData(phoneNumber, location, amount) {
    if (location.includes('Faso')) {
        location = 'Faso';
    }
    const url = `${BASE_URL}/rain`;

    const fetchOptions = {
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json"
        },
        method: 'POST',
        body: JSON.stringify({
            phoneNumber, location, amount
        })
    };

    const response = await fetch(url, fetchOptions);
    const data = await response.json();
    console.log(data);
    return data;
}
