package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *ApplicationConfigurationStatus) addCondition(ctype ApplicationConditionType, status corev1.ConditionStatus, reason, message string) {
	now := metav1.Now()
	c := &ApplicationCondition{
		Type:               ctype,
		LastUpdateTime:     now,
		LastTransitionTime: now,
		Status:             status,
		Reason:             reason,
		Message:            message,
	}
	m.Conditions = append(m.Conditions, *c)
}

// setConditionValue updates or creates a new condition
func (m *ApplicationConfigurationStatus) setConditionValue(ctype ApplicationConditionType, status corev1.ConditionStatus, reason, message string) {
	var c *ApplicationCondition
	for i := range m.Conditions {
		if m.Conditions[i].Type == ctype {
			c = &m.Conditions[i]
		}
	}
	if c == nil {
		m.addCondition(ctype, status, reason, message)
	} else {
		// check message ?
		if c.Status == status && c.Reason == reason && c.Message == message {
			return
		}
		now := metav1.Now()
		c.LastUpdateTime = now
		if c.Status != status {
			c.LastTransitionTime = now
		}
		c.Status = status
		c.Reason = reason
		c.Message = message
	}
}

// RemoveCondition removes the condition with the provided type.
func (m *ApplicationConfigurationStatus) RemoveCondition(ctype ApplicationConditionType) {
	for i, c := range m.Conditions {
		if c.Type == ctype {
			m.Conditions[i] = m.Conditions[len(m.Conditions)-1]
			m.Conditions = m.Conditions[:len(m.Conditions)-1]
			break
		}
	}
}

// GetCondition get existing condition
func (m *ApplicationConfigurationStatus) GetCondition(ctype ApplicationConditionType) *ApplicationCondition {
	for i := range m.Conditions {
		if m.Conditions[i].Type == ctype {
			return &m.Conditions[i]
		}
	}
	return nil
}

// IsConditionTrue - if condition is true
func (m *ApplicationConfigurationStatus) IsConditionTrue(ctype ApplicationConditionType) bool {
	if c := m.GetCondition(ctype); c != nil {
		return c.Status == corev1.ConditionTrue
	}
	return false
}

// IsReady returns true if ready condition is set
func (m *ApplicationConfigurationStatus) IsReady() bool { return m.IsConditionTrue(Ready) }

// IsNotReady returns true if ready condition is set
func (m *ApplicationConfigurationStatus) IsNotReady() bool { return !m.IsConditionTrue(Ready) }

// ConditionReason - return condition reason
func (m *ApplicationConfigurationStatus) ConditionReason(ctype ApplicationConditionType) string {
	if c := m.GetCondition(ctype); c != nil {
		return c.Reason
	}
	return ""
}

// Ready - shortcut to set ready contition to true
func (m *ApplicationConfigurationStatus) Ready(reason, message string) {
	m.SetConditionTrue(Ready, reason, message)
}

// NotReady - shortcut to set ready contition to false
func (m *ApplicationConfigurationStatus) NotReady(reason, message string) {
	m.SetConditionFalse(Ready, reason, message)
}

// SetError - shortcut to set error condition
func (m *ApplicationConfigurationStatus) SetError(reason, message string) {
	m.SetConditionTrue(Error, reason, message)
}

// ClearError - shortcut to set error condition
func (m *ApplicationConfigurationStatus) ClearError() {
	m.SetConditionFalse(Error, "NoError", "No error seen")
}

// SetConditionFalse updates or creates a new condition
func (m *ApplicationConfigurationStatus) SetConditionFalse(ctype ApplicationConditionType, reason, message string) {
	m.setConditionValue(ctype, corev1.ConditionFalse, reason, message)
}

// SetConditionTrue updates or creates a new condition
func (m *ApplicationConfigurationStatus) SetConditionTrue(ctype ApplicationConditionType, reason, message string) {
	m.setConditionValue(ctype, corev1.ConditionTrue, reason, message)
}

// RemoveAllConditions updates or creates a new condition
func (m *ApplicationConfigurationStatus) RemoveAllConditions() {
	m.Conditions = []ApplicationCondition{}
}

// ClearAllConditions updates or creates a new condition
func (m *ApplicationConfigurationStatus) ClearAllConditions() {
	for i := range m.Conditions {
		m.Conditions[i].Status = corev1.ConditionFalse
	}
}
