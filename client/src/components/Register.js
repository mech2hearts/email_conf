import React, { Component } from 'react'
import axios from 'axios'


export class Register extends Component{
  state = {
    email: '',
    message: ''
  }

  handleChange = (e) => {
    this.setState({
      [e.target.id] : e.target.value
    })
  }


  handleSubmit = (e) => {
    e.preventDefault()
    if(this.state.email!==""){
      const userInfo = {
        email: this.state.email
      }
      axios.post(`/api/`, userInfo)
        .then(res => {
          this.setState({
            message : res.data.message
          })
        })
    }
  }



  render(){
    return(
      <div className="container">
        <h3>Email Verification</h3>
        <div className="row z-depth-2">
          <div className="col s12 white z-depth-1">
            <form>
              <div className="input-field">
                <input type="text" id="email" value={this.state.email} onChange={this.handleChange} />
                <label htmlFor="email">Email</label>
              </div>
            </form>

            <div className="btn blue waves-effect waves-light" onClick={this.handleSubmit} >Login</div>
            <br />
            {
               (this.state.message === '') ? <br /> : <div>{this.state.message}</div>
            }
          </div>
        </div>
        <h4>About This Program</h4>
        <div className="row z-depth-2">
          <div className="col s12 white z-depth-1" style={{overflowY:"auto", height: 300}} >
            <b>Description:</b>
            <p>This program simulates the process of verifying newly registered accounts through
            the use of emails. <i>Note: Emails using Microsoft Outlook may alter the confirmation url.</i></p>
            <b>Components:</b>
            <p>
              ReactJS, MaterializeCSS, GO, MongoDB(mLab)
            </p>
            <b>Process:</b>
            <ol>
              <li>
                The user submits their email address into the form above.
              </li>
              <li>
                The email is sent to the back-end through an HTTP request using Axios.
              </li>
              <li>
                The GO back-end takes the email and checks to see if the email exists in the MongoDB database.
                <ul>
                  <li>
                    If the email does not exist in the database, the user will be sent a verification email.
                  </li>
                  <li>
                    If the email already exists, the user will be notified that an email has already been sent.
                  </li>
                </ul>
              </li>
              <li>
                The user will receive an email with the verification link. Upon submission of the link
                the user's email will then be verified and completely registered into the database.
              </li>
            </ol>
          </div>
        </div>
      </div>
    )
  }
}

export default Register
