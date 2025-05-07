document
  .getElementById('check-availability-button')
  .addEventListener('click', async () => {
    let html = `
      <form id="check-availability-form" action="/reservation" method="POST" novalidate class="needs validation">
        <div class="row g-3" id="reservation-dates-modal">
          <div class="col">
            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
          </div>
          <div class="col">
            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
          </div>
        </div>
      </form>
    `;
    const myDates = await attention.custom({
      title: "Check Bungalow's Availability",
      msg: html,
      preConfirm: () => {
        return [
          document.getElementById('start').value,
          document.getElementById('end').value,
        ];
      },
      willOpen: () => {
        const element = document.getElementById('reservation-dates-modal');
        const rp = new DateRangePicker(element, {
          format: 'yyyy-mm-dd',
          minDate: new Date(),
          showOnFocus: true,
        });
      },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled');
        document.getElementById('end').removeAttribute('disabled');
      },

      callback: (result) => {
        console.log('bungalow.js', result);
        let formEl = document.getElementById('check-availability-form');
        let formData = new FormData(formEl); // get form data of the formEl
        formData.append('csrf_token', csrfToken);
        const fetchResult = fetch('/reservation-json', {
          method: 'POST',
          body: formData,
        })
          .then((response) => response.json())
          .then((data) => {
            console.log('data', data);
            console.log(data.ok);
            console.log(data.message);
            // wenn ich hier returne, verlässt das nicht nur den then Block, sondern das ganze callback ???
            return { result, data };
          });
        // console.log('fetchResult', fetchResult, 'result', result);
        // // auch so bekomme ich die Daten nicht zurück ... - obwohl sie da sind
        return {
          result,
          fetchResult,
        };
      },
    });
    console.log('bungalow.js,myDates=', myDates);
  });
