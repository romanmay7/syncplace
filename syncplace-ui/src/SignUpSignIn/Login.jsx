import React,{useState, useEffect, useContext} from "react";
import {Link, useNavigate} from "react-router-dom";
import styled from "styled-components";
import Logo from "../assets/SPLogo.png";
import {ToastContainer,toast} from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import axios from "axios";
import { loginRoute } from "../utils/APIRoutes";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext
 
function Login() {
  const navigate = useNavigate();
  const[values,setValues] = useState({
    username: "", 
    password: "",
  });
  const { login } = useContext(AuthContext); // Access the login function

  const handleSubmit = async (event)=> {
    event.preventDefault();
    console.log("Validating Request");
    if (handleValidation()) {
      console.log("Sending Request to: ",loginRoute);
      const {password, username} = values;
     try{
          const {data,status} = await axios.post(loginRoute, {
           username,
           password,
         });

         if(status===200 || status===201 )
           {
             console.log("Login is Successfull: ",data.userName,);
             console.log("token: ",data.token,);
             //localStorage.setItem('token', data.token);
             //localStorage.setItem('syncplace-app-user',data.userName);
             login(data.userName, data.token); // Call the login function from AuthContext

             navigate("/joinroom");
           }
         else {
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
    const {username, password} = values;
    if (password === "") 
    {
         toast.error("Email and Password are required!", toastOptions);
         return false;
    } else if (username === "") 
    {
         toast.error("Email and Password are required", toastOptions);
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
              min="3"
            />
            <input 
              type="password"
              placeholder="Password"
              name="password"
              onChange={(e) => handleChange(e)}
            />
            <button type="submit">Login</button>
            <span >
                Don't have an Account yet ? <Link to = "/signup">Register</Link>
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
  gap: 1.3rem;
  align-items: center;
  background-color: 131324;
  margin: 0 auto; /* Center the container horizontally */

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
    width: 100%; /* Make form width responsive */
    max-width: 400px; /* Limit form width */

    input {
      background-color: #ffffff;
      padding: 0.5rem;
      border: 0.1rem solid #000000;
      border-radius: 0.4rem;
      color: black;
      width: 100%;
      font-size: 1rem;
      &:focus {
        border: #00000076 0.1rem solid #997af0;
        outline: none;
      }
    }
    button {
      background-color: #e6d3b1;
      color: black;
      padding: 0.7rem 1.5rem; 
      border: none;
      font-weight: bold;
      cursor: pointer;
      border-radius: 0.4rem;
      font-size: 1.2rem;
      text-transform: uppercase;
      transition: 0.5s ease-in-out;
      &:hover {
        background-color: #c4ad84;
      }
    }
  }
`;


export default Login;