const form = document.querySelector('.form-container')
const helperError = document.querySelector('.helper.error')
const helperSuccess = document.querySelector('.helper.success')

const longUrlInput = document.querySelector('#long-url')
const labelLongUrlInput = document.querySelector('.label[for="short-url"]')
const shortUrlInput = document.querySelector('#short-url')
const customId = document.querySelector('#custom-id')
const shortLink = document.querySelector('.short-link')

form.onsubmit = reduceLink

// helperError.style.display = 'block'

longUrlInput.value = ''
shortUrlInput.value = ''

switchInputId()
function switchInputId() {
  helperError.style = 'none'
  shortUrlInput.disabled = !customId.checked
  shortUrlInput.required = customId.checked
  const classAction = !customId.checked ? 'add' : 'remove'
  labelLongUrlInput.classList[classAction]('disabled')
}

async function reduceLink(e) {
  e.preventDefault()
  helperError.style = 'none'
  helperSuccess.style = 'none'
  try {
    const response = await fetch('http://localhost:4000/api', {
      method: 'POST',
      body: JSON.stringify({
        Url: longUrlInput.value,
        ShortUrl: shortUrlInput.value
      })
    });

    if (response.status === 500) {
      throw new Error(response)
    }
    const data = await response.json()
    console.log(response)
    console.log(data)

    shortUrlInput.value = ''
    longUrlInput.value = ''
    helperSuccess.style.display = 'block'
    console.log(data)
    const shortLinkValue =  `http://localhost:4000/api/redirect/${data.shortUrl}`
    shortLink.innerText = shortLinkValue
    shortLink.href = shortLinkValue
  } catch (error) {
    console.log(error)
    helperError.style.display = 'block' 
  }
}
