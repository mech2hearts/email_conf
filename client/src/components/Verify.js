import React, { Component } from 'react'
import axios from 'axios'
import { Redirect } from 'react-router-dom'


export class Verify extends Component{
  state = {
    message: '',
    redirect: false
  }


  componentDidMount () {
    const { id } = this.props.match.params
    console.log(id)
    axios.get(`/api/verify/${id}`)
      .then(res => {
        this.setState({
          message : res.data.message+" You will automatically be redirected to the main page."
        })
      })
      setTimeout(() => {
          this.setState({
              redirect: true,
          })
      }, 8000)
  }



  render(){
    return(
      <div className="container">
        {this.state.message}
        {
           (this.state.redirect === false) ? <br /> : <Redirect to='/' />
        }
      </div>
    )
  }
}

export default Verify
