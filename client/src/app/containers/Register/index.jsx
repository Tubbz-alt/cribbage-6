import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useInjectReducer, useInjectSaga } from 'redux-injectors';
import { useHistory } from 'react-router-dom';
import { authSaga } from '../../../auth/saga';
import { sliceKey, reducer, actions } from '../../../auth/slice';
import {
  Button,
  Container,
  Link,
  TextField,
  CssBaseline,
  Typography,
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  title: {
    fontSize: '2rem',
  },
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

const RegisterForm = () => {
  // hooks
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: authSaga });
  const history = useHistory();
  const dispatch = useDispatch();
  const [formData, setFormData] = useState({ id: '', name: '' });

  // event handlers
  const onSubmitForm = event => {
    event.preventDefault();
    dispatch(actions.register(formData.id, formData.name, history));
  };
  const onInputChange = event =>
    setFormData({ ...formData, [event.target.name]: event.target.value });

  const classes = useStyles();

  return (
    <Container component='main' maxWidth='sm'>
      <CssBaseline />
      <div className={classes.paper}>
        <Typography component='h1' className={classes.title}>
          Welcome to Cribbage!
        </Typography>
        <p>Play cribbage against your friends online. Get started now!</p>
        <form className={classes.form} onSubmit={onSubmitForm}>
          <TextField
            variant='outlined'
            margin='normal'
            required
            fullWidth
            label='Username'
            name='id'
            autoFocus
            onChange={onInputChange}
          />
          <TextField
            variant='outlined'
            margin='normal'
            required
            fullWidth
            name='name'
            label='Display Name'
            onChange={onInputChange}
          />
          <Button
            type='submit'
            fullWidth
            variant='contained'
            color='primary'
            className={classes.submit}
          >
            Register
          </Button>
        </form>
        <Link href='/' variant='body2'>
          Already have an account? Login here
        </Link>
      </div>
    </Container>
    // <div className='max-w-sm m-auto mt-12'>
    //     <input
    //       name='id'
    //       onChange={onInputChange}
    //       value={formData.id}
    //       placeholder='Username'
    //       required
    //       className='mt-2 form-input'
    //     ></input>
    //     <input
    //       name='name'
    //       onChange={onInputChange}
    //       value={formData.name}
    //       placeholder='Display name'
    //       required
    //       className='mt-2 form-input'
    //     ></input>
    //     <p className='mt-1 text-xs text-gray-600'>
    //       Already have an account?{' '}
    //       <span>
    //         <Link to='/' className='hover:text-gray-500 hover:underline'>
    //           Log in here.
    //         </Link>
    //       </span>
    //     </p>
    //     <input
    //       type='submit'
    //       value='register'
    //       className='mt-1 btn btn-primary'
    //     ></input>
    //   </form>
    // </div>
  );
};

export default RegisterForm;
