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
        ShortUrl: customId.checked ? shortUrlInput.value : ''
      })
    });

    const data = await response.json()
    if (response.status === 500) {
      throw new Error(data.message)
    }

    shortUrlInput.value = ''
    longUrlInput.value = ''
    helperSuccess.style.display = 'block'
    console.log(data)
    const shortLinkValue =  `http://localhost:4000/api/redirect/${data.shortUrl}`
    shortLink.innerText = shortLinkValue
    shortLink.href = shortLinkValue
    getAllUrl()
  } catch (error) {
    helperError.innerText = error.message
    helperError.style.display = 'block'
  }
}


const urlsArray = document.querySelector('.urls-array')
const urlsArrayCount = document.querySelector('.urls-array-count')

async function deleteUrl (ShortUrl) {
  helperError.style.display = 'none'
  helperSuccess.style.display = 'none'
  try {
    const response = await fetch(`http://localhost:4000/api/${ShortUrl}`, {
      method: 'DELETE'
    })

    const data = response.json()
    if (response.status === 500) {
      throw new Error(data.message)
    }

    getAllUrl()
  } catch (e) {
    console.log(e)
  }
}


const updateUrlInput = document.querySelector('#update-url')
async function updateUrl (ShortUrl) {
  helperError.style.display = 'none'
  helperSuccess.style.display = 'none'
  try {
    const response = await fetch(`http://localhost:4000/api/`, {
      method: 'PUT',
      body: JSON.stringify({
        Url: updateUrlInput.value,
        ShortUrl
      })
    })

    const data = response.json()
    if (response.status === 500) {
      throw new Error(data.message)
    }

    getAllUrl()
  } catch (e) {
    console.log(e)
  }
}

getAllUrl()
async function getAllUrl () {
  while (urlsArray.firstChild) {
    urlsArray.removeChild(urlsArray.firstChild);
  }

  try {
    const response = await fetch('http://localhost:4000/api')
    const data = await response.json()
    urlsArrayCount.innerText = `(${data ? data.length : '0'})`

    if (data === null) {
      throw new Error('No Data')
    }

    data.forEach((link) => {
      linkValue = `http://localhost:4000/api/redirect/${link.shortUrl}`
      let div = document.createElement('div')

      let linkContainer = document.createElement('p')

      div.classList.add('url-link')
      let a = document.createElement('a')
      a.classList.add('short-link')
      a.href = linkValue
      a.target = '_blank'
      a.innerText = linkValue
      linkContainer.appendChild(a)

      // Create link
      let countLink = document.createElement('span')
      countLink.textContent = `(${link.count})`
      countLink.classList.add('count-link')
      linkContainer.appendChild(countLink)

      div.appendChild(linkContainer)

      // add delete button
      let deleteButton = document.createElement('button')
      deleteButton.classList.add('link-delete')
      deleteButton.innerText = 'Supprimer'
      deleteButton.addEventListener('click', () => deleteUrl(link.shortUrl))

      let updateButton = document.createElement('button')
      updateButton.classList.add('link-delete')
      updateButton.innerText = 'Modifier'
      updateButton.addEventListener('click', () => updateUrl(link.shortUrl))

      let containerButton = document.createElement('div')
      containerButton.appendChild(deleteButton)
      containerButton.appendChild(updateButton)
      div.appendChild(containerButton)

      // add row
      urlsArray.appendChild(div)
    
    })
  } catch (e) {
    console.log(e)
  }
}

