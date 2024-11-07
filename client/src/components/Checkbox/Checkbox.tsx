import React from 'react';
import './Checkbox.scss';

interface CheckboxProps {
    checked: boolean;
    onChange: () => void;
}

const Checkbox: React.FC<CheckboxProps> = ({ checked, onChange }) => {
    return (
            <input
                type="checkbox"
                checked={checked}
                onChange={onChange}
                className="checkbox"
            />
    );
};

export default Checkbox;
