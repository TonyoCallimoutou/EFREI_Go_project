const form = document.querySelector('.form-container')
const helperError = document.querySelector('.helper.error')
const helperSuccess = document.querySelector('.helper.success')

const longUrlInput = document.querySelector('#long-url')
const labelLongUrlInput = document.querySelector('.label[for="short-url"]')
const shortUrlInput = document.querySelector('#short-url')
const customId = document.querySelector('#custom-id')

form.onsubmit = reduceLink

// helperError.style.display = 'block'

switchInputId()
function switchInputId() {
  shortUrlInput.disabled = !customId.checked
  shortUrlInput.required = customId.checked
  const classAction = !customId.checked ? 'add' : 'remove'
  labelLongUrlInput.classList[classAction]('disabled')
}

function reduceLink(e) {
  e.preventDefault()
  console.log(longUrlInput.value)
  console.log(shortUrlInput.value)
  try {
    fetch('http://localhost:4000', {
      method: 'POST',
      body: new FormData(document.querySelector('form')) // Collect form data
    })
      .then(response => response.text()) // Read response as text
      .then(data => alert(data)); // Alert the response
  } catch {
    helperError.style.display = 'block'
  }
}
