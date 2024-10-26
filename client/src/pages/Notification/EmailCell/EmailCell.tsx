import './EmailCell.scss';

interface EmailColumnProps {
    email: string;
}

const EmailColumn = ({email}: EmailColumnProps) => {
    return <div className='email-cell'>{email}</div>;
};

export default EmailColumn;