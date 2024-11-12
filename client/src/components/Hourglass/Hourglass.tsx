import React from 'react';
import './Hourglass.scss';

const Hourglass: React.FC = () => (
    <div className="hourglass">
        <div className="hourglass__sand-top"></div>
        <div className="hourglass__sand-bottom"></div>
        <div className="hourglass__frame"></div>
    </div>
);

export default Hourglass;
