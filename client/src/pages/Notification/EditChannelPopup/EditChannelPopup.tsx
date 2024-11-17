export interface EditChannelPopupProps {
  isDisplayed: boolean;
  setIsDisplayed: (arg: boolean) => void;
  onSubmit: () => void;
}

export interface EditChannelPopup {
  popup: ({
    isDisplayed,
    setIsDisplayed,
  }: EditChannelPopupProps) => React.ReactNode;
}
