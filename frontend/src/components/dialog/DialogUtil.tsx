export const RowContainer = ({
  children
}: {
  children: React.ReactNode;
}) => {
  return (
      <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
        {children}
      </div>
  );
};

export const LabelInputContainer = ({
  children
}: {
  children: React.ReactNode;
}) => {
  return (
      <div className="flex w-full flex-col space-y-2">
        {children}
      </div>
  );
};
