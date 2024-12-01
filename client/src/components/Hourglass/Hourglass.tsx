import React from 'react';
import './Hourglass.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

const Hourglass: React.FC = () => (
  <div className="hourglass">
    <div className="hourglass__frame-top">
      <SVGIcon iconName="sand-clock-top"/>
    </div>
    <div className="hourglass__frame-middle">
      <SVGIcon iconName="sand-clock-middle"/>
     </div>
     <div className="hourglass__frame-bottom">
       <SVGIcon iconName="sand-clock-bottom"/>
     </div>
   </div>
);

export default Hourglass;
