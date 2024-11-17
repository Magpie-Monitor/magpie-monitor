import './LabelInput.scss';

interface LabelInputParams {
  label: string;
  onChange?: (value: string) => void;
  placeholder?: string;
  value: string;
  disabled?: boolean;
  validationMessage?: (value: string) => string | null;
}

export const nonEmptyFieldValidation = (value: string) => {
  if (value.length === 0) {
    return 'Field cannot be empty';
  }

  return null;
};

export const fieldLengthValidation =
  (minLength: number, maxLegth: number) => (value: string) => {
    if (value.length < minLength || value.length > maxLegth) {
      return `Field length must be between ${minLength} and ${maxLegth}`;
    }

    return null;
  };

export const fieldPrefixValidation = (prefix: string) => (value: string) => {
  if (!value.startsWith(prefix)) {
    return `Field must start with ${prefix}`;
  }

  return null;
};

export const combineValidators =
  (validation: Array<(value: string) => string | null>) => (value: string) => {
    validation.reduce((accum, curr) => {
      const message = curr(value);
      if (message) {
        accum += message + '\n';
      }
      return accum;
    }, '');
  };

const LabelInput = ({
  label,
  onChange,
  placeholder,
  value,
  disabled,
  validationMessage,
}: LabelInputParams) => {
  return (
    <div className="label-input">
      <label className="label-input__label">{label}</label>
      <input
        type="text"
        disabled={disabled}
        value={value}
        placeholder={placeholder}
        onChange={(e) => (onChange ? onChange(e.target.value) : () => { })}
        className="label-input__input"
      />
      {validationMessage && validationMessage(value) && (
        <label className="label-input__validation">
          {validationMessage(value)}
        </label>
      )}
    </div>
  );
};

export default LabelInput;
