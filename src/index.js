import React from 'react'
import { render } from 'react-dom'

const compare = (a, b) => {
  const aFamily = a[1].params.device.familyType.toLowerCase()
  const aUUID = a[1].params.device.uuid.toLowerCase()
  const bFamily = b[1].params.device.familyType.toLowerCase()
  const bUUID = b[1].params.device.uuid.toLowerCase()
  const x = `${aFamily}-${aUUID}`
  const y = `${bFamily}-${bUUID}`
  if (x < y) {
    return -1
  }
  if (x > y) {
    return 1
  }
  return 0
}

class App extends React.Component {
  state = {
    data: {},
    filter: ''
  }

  async componentDidMount() {
    window.setInterval(() => {
      this.fetch()
    }, 1000)
  }

  fetch = async () => {
    const res = await fetch('/json')
    const data = await res.json()
    this.setState({
      data
    })
  }

  serviceHTTP = services => {
    return services && services.filter(s => s.type === 'http').length
  }

  onChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    })
  }

  render() {
    const entries = Object.entries(this.state.data)
      .filter(entry => {
        const params = JSON.stringify(entry[1].params).toLowerCase()
        return params.indexOf(this.state.filter.toLowerCase()) !== -1
      })
      .sort(compare)
    return (
      <div>
        <label htmlFor="filter">Filter:</label>
        <input id="filter" type="text" name="filter" onChange={this.onChange} />
        <table>
          <thead>
            <tr>
              <th />
              <th>UUID</th>
              <th>Family Type</th>
              <th>Label</th>
              <th>Name</th>
              <th>Firmware Version</th>
              <th>Website</th>
            </tr>
          </thead>
          <tbody>
            {entries.map((entry, i) => {
              const ip = entry[0]
              const params = entry[1].params
              return (
                <tr key={i}>
                  <td className="index">#{i + 1}</td>
                  <td>{params.device.uuid}</td>
                  <td>{params.device.familyType}</td>
                  <td>{params.device.label}</td>
                  <td>{params.device.name}</td>
                  <td>{params.device.firmwareVersion}</td>
                  <td>
                    {this.serviceHTTP(params.services) ? (
                      <a href={`http://${ip}`}>{ip}</a>
                    ) : null}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    )
  }
}

render(<App />, document.getElementById('react'))
