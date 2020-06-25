package sarama

type JoinGroupResponse struct {
	Version       int16
	ThrottleTime  int32
	Err           KError
	GenerationID  int32
	GroupProtocol string
	LeaderID      string
	MemberID      string
	Members       map[string][]byte
}

func (r *JoinGroupResponse) GetMembers() (map[string]ConsumerGroupMemberMetadata, error) {
	members := make(map[string]ConsumerGroupMemberMetadata, len(r.Members))
	for id, bin := range r.Members {
		meta := new(ConsumerGroupMemberMetadata)
		if err := decode(bin, meta); err != nil {
			return nil, err
		}
		members[id] = *meta
	}
	return members, nil
}

func (r *JoinGroupResponse) encode(pe packetEncoder) error {
	if r.Version >= 2 {
		pe.putInt32(r.ThrottleTime)
	}
	pe.putInt16(int16(r.Err))
	pe.putInt32(r.GenerationID)

	if err := pe.putString(r.GroupProtocol); err != nil {
		return err
	}
	if err := pe.putString(r.LeaderID); err != nil {
		return err
	}
	if err := pe.putString(r.MemberID); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.Members)); err != nil {
		return err
	}

	for memberID, memberMetadata := range r.Members {
		if err := pe.putString(memberID); err != nil {
			return err
		}

		if err := pe.putBytes(memberMetadata); err != nil {
			return err
		}
	}

	return nil
}

func (r *JoinGroupResponse) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version

	if version >= 2 {
		if r.ThrottleTime, err = pd.getInt32(); err != nil {
			return
		}
	}

	kerr, err := pd.getInt16()
	if err != nil {
		return err
	}

	r.Err = KError(kerr)

	if r.GenerationID, err = pd.getInt32(); err != nil {
		return
	}

	if r.GroupProtocol, err = pd.getString(); err != nil {
		return
	}

	if r.LeaderID, err = pd.getString(); err != nil {
		return
	}

	if r.MemberID, err = pd.getString(); err != nil {
		return
	}

	n, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	r.Members = make(map[string][]byte)
	for i := 0; i < n; i++ {
		memberID, err := pd.getString()
		if err != nil {
			return err
		}

		memberMetadata, err := pd.getBytes()
		if err != nil {
			return err
		}

		r.Members[memberID] = memberMetadata
	}

	return nil
}

func (r *JoinGroupResponse) key() int16 {
	return 11
}

func (r *JoinGroupResponse) version() int16 {
	return r.Version
}

func (r *JoinGroupResponse) headerVersion() int16 {
	return 0
}

func (r *JoinGroupResponse) requiredVersion() KafkaVersion {
	switch r.Version {
	case 2:
		return V0_11_0_0
	case 1:
		return V0_10_1_0
	default:
		return V0_9_0_0
	}
}
