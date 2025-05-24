const validateContact = (rule, value, callback) => {
  const email = form.value.email?.trim()
  const phone = form.value.phone?.trim()
  if (!email && !phone) {
    // Please enter at least an email or a phone number
    callback(new Error('Either email or phone must be provided'))
  } else {
    callback()
  }
}

const validateLabels = (rule, value, callback) => {
  if (!Array.isArray(value)) {
    callback();
    return;
  }

  const tooLong = value.find(label => label.length > 32);
  if (tooLong) {
    callback(new Error(`Label ${tooLong} exceeds 32 characters`));
    return;
  }

  if (value.length > 16) {
    callback(new Error('You can select up to 16 labels only'));
    return;
  }

  callback();
}

export const validateAccount = {
  firstname: [{ required: true, min: 2, max: 32, message: 'Please enter firstname', trigger: 'blur' }],
  lastname: [{ required: true, min: 2, max: 32, message: 'Please enter lastname', trigger: 'blur' }],
  email: [
    { validator: validateContact, trigger: 'blur' },
    { type: 'email', message: 'invalid email format', trigger: 'blur' },
    { min: 5, max: 64, message: 'Email too long', trigger: 'blur' },
  ],
  phone: [{ min: 6, max: 20, validator: validateContact, trigger: 'blur' }],
  password: [{
    required: true, message: 'Please enter password',
    trigger: 'blur', min: 8, max: 32,
    pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,32}$/,
    message: 'Password must be 8-32 chars with at least one uppercase, lowercase and number',
  }],
  level: [{ required: true, message: 'Please enter level', trigger: 'change' }],
  labels: [{ validator: validateLabels, trigger: 'change' }],
}
