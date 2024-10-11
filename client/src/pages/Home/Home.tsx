import Table, { TableColumn } from 'components/Table/Table';

interface UserRow {
  name: string;
  surname: string;
  phoneNumber: number;
  [key: string]: string | number;
}

const customTextField = (row: UserRow): React.ReactNode => {
  return (
    <div> Customdsaksjdsdasfkjsdfksjldjldjasdaslkhdasjlhdasjhdlkasjdkjaajshdlak {row.name}</div>
  );
};

const TableData: {
  columns: TableColumn<UserRow>[];
  rows: UserRow[];
} = {
  columns: [
    {
      header: 'Name',
      columnKey: 'name',
      // customComponent: customTextField,
    },
    {
      header: 'Surname',
      columnKey: 'surname',
      customComponent: customTextField,
    },
    {
      header: 'Phone Number',
      columnKey: 'phoneNumber',
      customComponent: customTextField,
    },
  ],
  rows: [
    {
      name: 'John',
      surname: 'One',
      phoneNumber: 132,
    },
    {
      name: 'John',
      surname: 'Two',
      phoneNumber: 132,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'One',
      phoneNumber: 132,
    },
    {
      name: 'John',
      surname: 'Two',
      phoneNumber: 132,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
    {
      name: 'John',
      surname: 'Three',
      phoneNumber: 152,
    },
  ],
};

const Home = () => {
  return (
    <div style={{ width: '80vw', height: '50vh', backgroundColor: '#0C1926' }}>
      <Table rows={TableData.rows} columns={TableData.columns} />
    </div>
  );
};

export default Home;
