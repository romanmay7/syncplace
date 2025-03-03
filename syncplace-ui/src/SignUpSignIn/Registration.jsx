import React,{useState, useEffect} from "react";
import {Link, useNavigate} from "react-router-dom";
import styled from "styled-components";
import Logo from "../assets/SPLogo.png";
import {ToastContainer,toast} from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import axios from "axios";
import { registerRoute } from "../utils/APIRoutes";
 
function Register() {
  const navigate = useNavigate();
  const[values,setValues] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const handleSubmit = async (event)=> {
    event.preventDefault();
    console.log("Validating Request");
    if (handleValidation()) {
      console.log("Sending Request to: ",registerRoute);
      const {password, username, email} = values;
     try{  
        const {data,status} = await axios.post(registerRoute, {
          username,
          email,
          password,
        });

        if(status===200 || status===201 ) {
            //localStorage.setItem('syncplace-app-user',JSON.stringify(data.userName));
            console.log("New user created: ",data.userName);
            alert("New user created: ",data.userName);
            navigate("/login");
         }
         else 
         {
             toast.error(data.error, toastOptions);
         }  
       } catch(error) 
         {
          if (error.response && error.response.data && error.response.data.error)
             {
              // Access the error message from the server's response
              toast.error(error.response.data.error, toastOptions);
             } else
             {
               // Handle cases where the error response is not as expected
               toast.error("An unexpected error occurred.", toastOptions);
               console.error("Axios error:", error); // Log the full error for debugging
             }
         }
      }

  };

  const toastOptions = {
    position: "bottom-right",
    autoClose: 8000,
    pauseOnHover: true,
    draggable: true,
    theme: "dark",
  };

  const handleValidation = () => {
    const {username, email, password, confirmPassword} = values;
    if (password !== confirmPassword) 
    {
         toast.error("Password Confirmation has Failed!", toastOptions);
         return false;

    } else if (username.length < 3) 
    {
      toast.error("Username should be greater than 3 characters", toastOptions);
      return false;
    } else if(email.length < 3)
    {
        toast.error("Email should be greater than 3 characters", toastOptions);
        return false;
    }else if(password.length < 8)
    {
        toast.error("Password should be not shorter than 8 characters", toastOptions);
        return false;
    }

    return true;
  };

  const handleChange = (event) => {
    setValues({...values,[event.target.name]:event.target.value});

  }
    return (
      <>
        <FormContainer>
          <form onSubmit ={(event)=> handleSubmit(event)}>
            <div className ="brand">
              <img src={Logo} alt="" />
            </div>
            <input 
              type="text"
              placeholder="Username"
              name="username"
              onChange={(e) => handleChange(e)}
            />
            <input 
              type="email" 
              placeholder="Email"
              name="email"
              onChange={(e) => handleChange(e)}
            />
            <input 
              type="password"
              placeholder="Password"
              name="password"
              onChange={(e) => handleChange(e)}
            />
            <input 
              type="password"
              placeholder="Confirm Password"
              name="confirmPassword"
              onChange={(e) => handleChange(e)}
            />
            <button type="submit">Create User Account</button>
            <span>
                Already have an Account ? <Link to = "/login">Login</Link>
            </span>
          </form>
          <ToastContainer />
        </FormContainer>
      </>
    )
}

const FormContainer = styled.div`
height: 80vh;
 width: 90%; /* Changed to percentage for responsiveness */
 max-width: 170vh; /* Added max-width to limit size on larger screens */
 display: flex;
 flex-direction: column;
 justify-content: center;
 gap: 1rem;
 align-items: center;
 background-color: 131324;
 .brand {
    display: flex;
    align-items: center;
    gap: 1rem;
    justify-content: center;
    img {
      height: 5rem;
    }
      h1 {
        color: white;
        text-transform: uppercase;
      }
 }
form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  background-color: #ece6d8;
  border-radius: 2rem;
  padding: 3rem 5rem;
  input {
    background-color: #ffffff;
    padding: 1.2 rem;
    border:0.1rem solid #000000.;
    border-radius: 0%.4rem;
    color: black;
    width: 100%;
    font-size:1rem;
    &:focus {
       border: #00000076.1rem solid #997af0;
       outline: none;
    }
  }
  button {
    background-color: #e6d3b1;
    color: black;
    padding 1rem 2rem;
    border: none;
    font-weight: bold;
    cursor: pointer;
    border-radius: 0%.4rem;
    font-size: 1.1rem;
    text-transform: uppercase;
    transition: 0.5s ease-in-out;
    &:hover {
      background-color:#c4ad84
    }
  }      
}

`;


export default Register;